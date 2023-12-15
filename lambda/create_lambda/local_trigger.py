import sys
import json
import lambda_function as c

if __name__ == "__main__":
    # Access command line arguments
    args = sys.argv[1:]
    print(args)
    if not args:
        print("No arguments provided.")
        sys.exit(1)

    # Assuming the argument is a JSON-formatted string
    hello = dict()
    print("entered")
    try:
        # Parse the JSON string into a dictionary
        # data = json.loads(json_string)
        hello['service_id']=args
        response = c.lambda_handler(hello, "context")
        # # Accessing keys and values
        # for key, value in data.items():
        #     print(f"Key: {key}, Value: {value}")
        print(response.body.json.dumps())

    except json.JSONDecodeError as e:
        print(f"error decoding json: {e}")
        sys.exit(1)