import json
import sys
import requests

def main():
    try:
        input_data = json.load(sys.stdin)
        payload_url = input_data.get("payload_url")

        if not payload_url:
            sys.stderr.write("Error: payload_url not provided in input\n")
            sys.exit(1)

        try:
            response = requests.get(payload_url)
            response.raise_for_status()
            payload_data = response.json()
        except (requests.RequestException, ValueError) as e:
            sys.stderr.write(f"Error fetching or parsing payload from URL: {e}\n")
            sys.exit(1)

        if not isinstance(payload_data, list):
            sys.stderr.write("payload from URL must contain a list of connector configurations.\n")
            sys.exit(1)

        output_payload = {
            "connectors": json.dumps(payload_data)
        }

        json.dump(output_payload, sys.stdout)

    except Exception as e:
        sys.stderr.write(f"Error: {str(e)}\n")
        sys.exit(1)

if __name__ == "__main__":
    main()