import sys
import yaml
from dotenv import dotenv_values

def get_env_var(key):
    value = env_vars.get(key)
    if value is None:
        print(f"Error: Environment variable {key} is not set.")
        sys.exit(1)
    return value

def fill_template(template_path, output_path):
    # Load the YAML template
    with open(template_path, 'r') as f:
        template = yaml.safe_load(f)

    # Update the template with environment variables
    template['mongodb']['uri'] = get_env_var('MONGODB_URI')
    template['mongodb']['database'] = get_env_var('MONGODB_DATABASE')

    template['pocket']['rpc_url'] = get_env_var('POKT_RPC_URL')
    template['pocket']['chain_id'] = get_env_var('POKT_CHAIN_ID')
    template['pokt_signer']['private_key'] = get_env_var('POKT_PRIVATE_KEY')

    template['ethereum']['rpc_url'] = get_env_var('ETH_RPC_URL')
    template['ethereum']['chain_id'] = get_env_var('ETH_CHAIN_ID')
    template['wpokt_signer']['private_key'] = get_env_var('ETH_PRIVATE_KEY')

    # Save the filled template to the output file
    with open(output_path, 'w') as f:
        yaml.dump(template, f)

# Check if the correct number of arguments is provided
if len(sys.argv) != 4:
    print("Usage: python fill_template.py <env_file> <template_file> <output_file>")
else:
    # Get the input and output file paths from command-line arguments
    env_file = sys.argv[1]
    template_file = sys.argv[2]
    output_file = sys.argv[3]

    # Load the environment variables
    env_vars = dotenv_values(env_file)

    # Call the function to fill the template
    fill_template(template_file, output_file)
