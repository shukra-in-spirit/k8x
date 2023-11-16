import boto3
from boto3.dynamodb.conditions import Key
import pandas as pd
from numpy import array, asarray
from keras.models import Sequential
from keras.layers import LSTM, Dense, Bidirectional

# get the data from the dynamodb table for the service
def query_dynamodb_table(partition_key_value):
    dynamodb = boto3.resource('dynamodb')
    table = dynamodb.Table('b_k8x_t')

    # Define the query parameters
    key_condition_expression = Key('service_id').eq(partition_key_value)
    
    # Perform the query
    response = table.query(
        KeyConditionExpression=key_condition_expression
    )
    
    # Retrieve just the items
    items = response['Items']

    # pagination in dynamodb
    while 'LastEvaluatedKey' in response:
        # getting page n of items
        response = table.query(
            KeyConditionExpression=key_condition_expression,
            ExclusiveStartKey=response['LastEvaluatedKey']
        )
        # appending to items list
        items.append(response['Items'])
    #print(items)
    return items

# create the dataframes
def create_dataframes(cpu_memory_data):
    # normalise the json to dataframe
    cpu_memory_data_pd = pd.json_normalize(cpu_memory_data, meta=['timestamp', 'cpu', 'memory'])
    
    # create the cpu list
    cpu_df = cpu_memory_data_pd['cpu'].tolist()
    
    # create the memory list
    memory_df = cpu_memory_data_pd['memory'].tolist()

    # convert the decimal.Decimal to float32
    cpu_df = asarray(cpu_df).astype('float32')
    memory_df = asarray(memory_df).astype('float32')

    return cpu_df, memory_df

# split a univariate sequence
def split_sequence(sequence, n_steps):
    X, y = list(), list()
    for i in range(len(sequence)):
        # find the end of this pattern
        end_ix = i + n_steps
        # check if we are beyond the sequence
        if end_ix > len(sequence)-1:
            break
        # gather input and output parts of the pattern
        seq_x, seq_y = sequence[i:end_ix], sequence[end_ix]
        X.append(seq_x)
        y.append(seq_y)
    return array(X), array(y)

# create and fit the model
def create_and_fit_model(training_data, n_steps, n_features, hidden_layers, epochs):
    # split into X and y for fit
    X, y = split_sequence(training_data, n_steps)

    # reshape x
    X = X.reshape(X.shape[0], X.shape[1], n_features)

    # define model
    model = Sequential()

    # add hidden layers
    model.add(Bidirectional(LSTM(hidden_layers, activation='relu'), input_shape=(n_steps, n_features)))
    
    # add output layer
    model.add(Dense(1))

    # compile the model architecture
    model.compile(optimizer='adam', loss='mse')

    # fit model
    model.fit(X, y, epochs=epochs, verbose=0)

    return model

# upload the model to s3
def upload_model_to_s3(local_path, s3_path):
    s3 = boto3.client('s3')
    try:
        # Upload the file
        s3.upload_file(local_path, 'k8x', s3_path)
        print(f"File uploaded successfully to s3://k8x/{s3_path}")
    except Exception as e:
        print(f"Error uploading file to S3: {e}")
    return

def calculate_average(data):
    sum = 0
    for i in range(0,len(data)):
        sum += data[i]
    return sum/i

def lambda_handler(event, context):
   
    # get the training data from dynamodb
    training_data_dd = query_dynamodb_table(event['service_id'])

    # convert the data to pandas data frame
    cpu_training_df, memory_training_df = create_dataframes(training_data_dd)
    print("2")

    ## CPU ##
    # set cpu model training parameters
    cpu_n_steps = 1
    cpu_n_features = 1
    cpu_hidden_layers = 50
    cpu_epochs = 200

    # ceate the cpu model 
    cpu_model = create_and_fit_model(cpu_training_df, cpu_n_steps, cpu_n_features, cpu_hidden_layers, cpu_epochs)
    print('3')
    
    # save cpu model
    cpu_model_s3_path = event['service_id'] + '/cpu_model.h5'
    cpu_model_local_path = '~/tmp/' + cpu_model_s3_path
    cpu_model.save(cpu_model_local_path)

    # calculate cpu average
    cpu_avg = calculate_average(cpu_training_df)
    print('4')

    # upload model to s3
    upload_model_to_s3(cpu_model_local_path, cpu_model_s3_path)
    print('5')

    ## Memory ##
    # set memory model training parameters
    memory_n_steps = 1
    memory_n_features = 1
    memory_hidden_layers = 50
    memory_epochs = 200

    # ceate the memory model 
    memory_model = create_and_fit_model(memory_training_df, memory_n_steps, memory_n_features, memory_hidden_layers, memory_epochs)

    # save memory model
    memory_model_s3_path = event['service_id'] + '/memory_model.h5'
    memory_model_local_path = '~/tmp/' + memory_model_s3_path
    memory_model.save(memory_model_local_path)

    # calculate memory average
    memory_avg = calculate_average(memory_training_df)

    # upload model to s3
    upload_model_to_s3(memory_model_local_path, memory_model_s3_path)

    return {
        'statusCode': 200,
        'body': {
            "cpu": cpu_avg,
            "memory": memory_avg
        }
    }

hello = dict()
hello['service_id']="mrs-vas-1"
lambda_handler(hello,'jello')