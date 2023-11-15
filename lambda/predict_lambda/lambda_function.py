import boto3
from keras.models import load_model
import pandas as pd
from numpy import asarray

# get the model from the s3
def get_model_from_s3(path, local_name):
    # create the client for s3
    s3 = boto3.client('s3')
    
    # download the model and store it in the local path
    try:
        s3.download_file('k8x', path, local_name)
        print(f"File downloaded successfully")
    except Exception as e:
        print(f"Error downloading file from S3: {e}")
    
    # load the model using keras
    model = load_model(local_name)

    return model

# get the event data and convert it into a df list
def get_data_from_event(event):
    prediction_raw_data = event['history']
    # normalise the json to dataframe
    cpu_memory_data_pd = pd.json_normalize(prediction_raw_data, meta=['timestamp', 'cpu', 'memory'])
    
    # create the cpu list
    cpu_df = cpu_memory_data_pd['cpu'].tolist()
    
    # create the memory list
    memory_df = cpu_memory_data_pd['memory'].tolist()

    # convert the decimal.Decimal to float32
    cpu_df = asarray(cpu_df).astype('float32')
    memory_df = asarray(memory_df).astype('float32')

    return cpu_df, memory_df


def lambda_handler(event, context):
    # get the prediction data from event
    cpu_prediction_data, memory_prediction_data = get_data_from_event(event)
    
    ## CPU ##
    # get the model from s3
    cpu_model = get_model_from_s3('mrs-vas-1/cpu_model.h5', 'cpu_model.h5')
    
    # prepare for prediction
    n_steps, n_features, n_predictions = 1,1,1
    x_input = cpu_prediction_data.reshape(n_predictions, n_steps, n_features)
    
    # make the model predict
    cpu_value = cpu_model.predict(x_input, verbose=0)

    ## Memory ##
    # get the model from s3
    memory_model = get_model_from_s3('mrs-vas-1/memory_model.h5', 'memory_model.h5')
    
    # prepare for prediction
    n_steps, n_features, n_predictions = 1,1,1
    x_input = memory_prediction_data.reshape(n_predictions, n_steps, n_features)
    
    # make the model predict
    memory_value = memory_model.predict(x_input, verbose=0)
    # call service server 
    
    print("predicted_memory " + str(memory_value[0][0]))
    print("predicted_cpu " + str(cpu_value[0][0]))

    return {
        'statusCode': 200,
        'body': {
            "cpu": cpu_value[0][0],
            "memory": memory_value[0][0]
        }
    }

event = {
    "service_id": "mrs-vas-1",
    "history": [
        {
            "timestamp": "2023-10-07 12:40:00",
            "cpu": "0.00106",
            "memory": "40787968"
        }
    ],
    "hypertuning_parameters": {
        "epochs":"",
    }

}
lambda_handler(event,'hello')