package main

import (
	"fmt"
	"image"
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	marsToEarth := make(chan []Message)
	go earthReceiver(marsToEarth)

	gridSize := image.Point{X: 20, Y: 10}
	grid := NewMarsGrid(gridSize)
	rover := make([]*RoverDriver, 2)
	for i := range rover {
		rover[i] = startDriver(fmt.Sprint("Марсоход ", i), grid, marsToEarth)
	}
	time.Sleep(60 * time.Second)
}

// Message содержит сообщение, отправленное с Марса до Земли
type Message struct {
	Pos       image.Point
	LifeSigns int
	Rover     string
}

const (
	// Длина марсианского дня
	dayLength = 24 * time.Second
	// Продолжительность, во время которого
	// сообщения можно передать с марсохода до Земли
	receiveTimePerDay = 2 * time.Second
)

// earthReceiver получает сообщения, отправленные с Марса
// Так как связь ограничена, принимаются только сообщения
// для некоторого часа марсианского дня
func earthReceiver(msgc chan []Message) {
	for {
		time.Sleep(dayLength - receiveTimePerDay)
		receiveMarsMessages(msgc)
	}
}

// receiveMarsMessages получает сообщения, отправленные с Марса
// для данной продолжительности
func receiveMarsMessages(msgc chan []Message) {
	finished := time.After(receiveTimePerDay)
	for {
		select {
		case <-finished:
			return
		case ms := <-msgc:
			for _, m := range ms {
				log.Printf("земля получает доклад об уровне жизни %d из %s в %v", m.LifeSigns, m.Rover, m.Pos)
			}
		}
	}
}

func startDriver(name string, grid *MarsGrid, marsToEarth chan []Message) *RoverDriver {
	var o *Occupier
	// Попытка получить случайную точку продолжается до тех пор, пока не будет найдена та,
	// что сейчас не занята
	for o == nil {
		startPoint := image.Point{X: rand.Intn(grid.Size().X), Y: rand.Intn(grid.Size().Y)}
		o = grid.Occupy(startPoint)
	}
	return NewRoverDriver(name, o, marsToEarth)
}

// Структура Radio представляет радио передатчик, что может отправить
// сообщение на Землю
type Radio struct {
	fromRover chan Message
}

// SendToEarth отправляет сообщение на Землю.
// Успешность видна сразу - реальное сообщение
// можно скопировать в буфер и передать позже
func (r *Radio) SendToEarth(m Message) {
	r.fromRover <- m
}

// NewRadio возвращает новый экземпляр Radio, что
// отправляет сообщения на канал toEarth
func NewRadio(toEarth chan []Message) *Radio {
	r := &Radio{
		fromRover: make(chan Message),
	}
	go r.run(toEarth)
	return r
}

// запуск сообщений буфера, отправляемых марсоходом до дня, пока
// они не смогут отправиться на Землю
func (r *Radio) run(toEarth chan []Message) {
	var buffered []Message
	for {
		toEarth1 := toEarth
		if len(buffered) == 0 {
			toEarth1 = nil
		}
		select {
		case m := <-r.fromRover:
			buffered = append(buffered, m)
		case toEarth1 <- buffered:
			buffered = nil
		}
	}
}

// RoverDriver ведет марсоход по поверхности Марса
type RoverDriver struct {
	commandc chan command
	occupier *Occupier
	name     string
	radio    *Radio
}

// NewRoverDriver начинает новый RoverDriver и возвращает его
func NewRoverDriver(
	name string,
	occupier *Occupier,
	marsToEarth chan []Message,
) *RoverDriver {
	r := &RoverDriver{
		commandc: make(chan command),
		occupier: occupier,
		name:     name,
		radio:    NewRadio(marsToEarth),
	}
	go r.drive()
	return r
}

type command int

const (
	right command = 0
	left  command = 1
)

// drive отвечает за ведение марсохода. Ожидается,
// что он начнется в горутине
func (r *RoverDriver) drive() {
	log.Printf("%s начальная позиция %v", r.name, r.occupier.Pos())
	direction := image.Point{X: 1, Y: 0}
	updateInterval := 250 * time.Millisecond
	nextMove := time.After(updateInterval)
	for {
		select {
		case c := <-r.commandc:
			switch c {
			case right:
				direction = image.Point{
					X: -direction.Y,
					Y: direction.X,
				}
			case left:
				direction = image.Point{
					X: direction.Y,
					Y: -direction.X,
				}
			}
			log.Printf("%s новое направление %v", r.name, direction)
		case <-nextMove:
			nextMove = time.After(updateInterval)
			newPos := r.occupier.Pos().Add(direction)
			if r.occupier.MoveTo(newPos) {
				log.Printf("%s перемещение на %v", r.name, newPos)
				r.checkForLife()
				break
			}
			log.Printf("%s заблокирован при попытке перемещения из %v в %v", r.name, r.occupier.Pos(), newPos)
			// Случайно выбирается одно из других случайных направлений
			// Далее мы попробуем передвинуться в новое направление
			dir := rand.Intn(3) + 1
			for i := 0; i < dir; i++ {
				direction = image.Point{
					X: -direction.Y,
					Y: direction.X,
				}
			}
			log.Printf("%s новое случайное направление %v", r.name, direction)
		}
	}
}

func (r *RoverDriver) checkForLife() {
	// Успешное перемещение на новую позицию
	sensorData := r.occupier.Sense()
	if sensorData.LifeSigns < 900 {
		return
	}
	r.radio.SendToEarth(Message{
		Pos:       r.occupier.Pos(),
		LifeSigns: sensorData.LifeSigns,
		Rover:     r.name,
	})
}

// Left поворачивает марсоход налево (90° против часовой стрелки).
func (r *RoverDriver) Left() {
	r.commandc <- left
}

// Right поворачивает марсоход направо (90° по часовой стрелке).
func (r *RoverDriver) Right() {
	r.commandc <- right
}

// MarsGrid представляет сетку некоторых поверхностей
// Марса. Это может использоваться конкурентно несколькими горутинами
type MarsGrid struct {
	bounds image.Rectangle
	mu     sync.Mutex
	cells  [][]cell
}

// SensorData содержит информацию о том, что находится в данной точке сетки
type SensorData struct {
	LifeSigns int
}

type cell struct {
	groundData SensorData
	occupier   *Occupier
}

// NewMarsGrid возвращает новый MarsGrid указанного размера
func NewMarsGrid(size image.Point) *MarsGrid {
	grid := &MarsGrid{
		bounds: image.Rectangle{
			Max: size,
		},
		cells: make([][]cell, size.Y),
	}
	for y := range grid.cells {
		grid.cells[y] = make([]cell, size.X)
		for x := range grid.cells[y] {
			cell := &grid.cells[y][x]
			cell.groundData.LifeSigns = rand.Intn(1000)
		}
	}
	return grid
}

// Size возвращает Point, что представляет размер сетки
func (g *MarsGrid) Size() image.Point {
	return g.bounds.Max
}

// Occupy занимает клетку в данной точке в сетке.
// Возвращается nil, если точка уже занята, или если точка находится
// за пределами сетки. В противном случае значение можно использовать
// для перемещения в другое место сетки
func (g *MarsGrid) Occupy(p image.Point) *Occupier {
	g.mu.Lock()
	defer g.mu.Unlock()
	cell := g.cell(p)
	if cell == nil || cell.occupier != nil {
		return nil
	}
	cell.occupier = &Occupier{
		grid: g,
		pos:  p,
	}
	return cell.occupier
}

func (g *MarsGrid) cell(p image.Point) *cell {
	if !p.In(g.bounds) {
		return nil
	}
	return &g.cells[p.Y][p.X]
}

// Occupier представляет занятую клетку сетки
type Occupier struct {
	grid *MarsGrid
	pos  image.Point
}

// MoveTo передвигает занятую клетку на другую клетку сетки
// Сообщается, было ли перемещение успешным
// Может закончится неудачей, если была попытка перемещения за пределы
// сетки, или потому что пунктом назначения является занятая
// клетка. В случае неудачи занятая клетка не перемещается и остается на прежнем месте.
func (o *Occupier) MoveTo(p image.Point) bool {
	o.grid.mu.Lock()
	defer o.grid.mu.Unlock()
	newCell := o.grid.cell(p)
	if newCell == nil || newCell.occupier != nil {
		return false
	}
	o.grid.cell(o.pos).occupier = nil
	newCell.occupier = o
	o.pos = p
	return true
}

// Sense возвращает сенсорные данные из текущей клетки
func (o *Occupier) Sense() SensorData {
	o.grid.mu.Lock()
	defer o.grid.mu.Unlock()
	return o.grid.cell(o.pos).groundData
}

// Pos возвращает текущую позицию сетки занятой клетки
func (o *Occupier) Pos() image.Point {
	return o.pos
}
