import os
import argparse

def set_environment_variables(env):
    variables = {
        'dev': {
            'ENV': 'dev',
            'CONFIG_PATH': './config/dev.yaml',
            'DB_HOST': 'localhost',
            'DB_PORT': '5432',
            'DB_NAME': 'sso_dev',
            'DB_USER': 'dev_user',
            'DB_PASSWORD': 'dev_password',
            'REDIS_HOST': 'localhost',
            'REDIS_PORT': '6379',
            'JWT_SECRET': 'dev_secret',
        },
        'prod': {
            'ENV': 'prod',
            'CONFIG_PATH': './config/prod.yaml',
            'DB_HOST': 'prod_db_host',
            'DB_PORT': '5432',
            'DB_NAME': 'sso_prod',
            'DB_USER': 'prod_user',
            'DB_PASSWORD': 'prod_password',
            'REDIS_HOST': 'prod_redis_host',
            'REDIS_PORT': '6379',
            'JWT_SECRET': 'prod_secret',
        },
        'test': {
            'ENV': 'test',
            'CONFIG_PATH': './config/test.yaml',
            'DB_HOST': 'localhost',
            'DB_PORT': '5432',
            'DB_NAME': 'sso_test',
            'DB_USER': 'test_user',
            'DB_PASSWORD': 'test_password',
            'REDIS_HOST': 'localhost',
            'REDIS_PORT': '6379',
            'JWT_SECRET': 'test_secret',
        }
    }

    if env not in variables:
        print(f"Invalid environment: {env}")
        return

    for key, value in variables[env].items():
        os.environ[key] = value
        print(f"Set {key}={value}")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Set environment variables for SSO project")
    parser.add_argument("env", choices=["dev", "prod", "test"], help="Environment to set variables for")
    args = parser.parse_args()

    set_environment_variables(args.env)