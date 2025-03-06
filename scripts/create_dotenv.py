import os

def read_env_file(file_path):
    with open(file_path, 'r') as file:
        return file.readlines()

def merge_env_files(file_contents):
    merged_lines = []
    seen_keys = set()

    for lines in file_contents:
        for line in lines:
            line = line.strip()
            if line.startswith('#') or not line:
                if line not in merged_lines:
                    merged_lines.append(line)
            else:
                key = line.split('=')[0].strip()
                if key not in seen_keys:
                    merged_lines.append(line)
                    seen_keys.add(key)

    return merged_lines

def prompt_for_values(merged_lines):
    final_lines = []
    for line in merged_lines:
        if not line.startswith('#') and '=' in line:
            key, default_value = line.split('=', 1)
            value = input(f"Enter value for {key} (default: {default_value}): ").strip()
            if not value:
                value = default_value
            final_lines.append(f"{key}={value}")
        else:
            final_lines.append(line)
    return final_lines

def write_merged_env_file(output_path, final_lines):
    with open(output_path, 'w') as file:
        for line in final_lines:
            file.write(line + '\n')

def main():
    script_dir = os.path.dirname(os.path.abspath(__file__))
    project_root = os.path.dirname(script_dir)
    
    template_paths = [
        os.path.join(project_root, 'docs', '.env.example'),
        os.path.join(project_root, 'sso-sdk', 'docs', '.env.example')
    ]
    
    output_path = os.path.join(project_root, '.env')

    file_contents = []
    for template_path in template_paths:
        if os.path.exists(template_path):
            file_contents.append(read_env_file(template_path))
        else:
            print(f"Warning: Template file not found at {template_path}")

    if not file_contents:
        print("Error: No template files found.")
        return

    merged_lines = merge_env_files(file_contents)
    final_lines = prompt_for_values(merged_lines)
    write_merged_env_file(output_path, final_lines)

    print(f"Merged .env file created successfully at {output_path}")

if __name__ == "__main__":
    main()