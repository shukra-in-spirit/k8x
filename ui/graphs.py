import tkinter as tk
from datetime import datetime
import random

def plot_metrics_in_tkinter(data, canvas, fig, ax):
    # Extract data from the JSON list
    timestamps = [entry['timestamp'] for entry in data]
    metric1_values = [entry['metric1'] for entry in data]
    metric2_values = [entry['metric2'] for entry in data]

    ax.clear()

    # Convert timestamps to datetime objects
    timestamps = [datetime.fromisoformat(timestamp) for timestamp in timestamps]

    ax.plot(timestamps, metric1_values, label='actual', marker='o')
    ax.plot(timestamps, metric2_values, label='predicted', marker='o')

    # Configure the plot
    ax.set_xlabel('Timestamp')
    ax.set_ylabel('Values')
    ax.set_title('Comparison of actual and predicted over 1 day')
    ax.legend()
    #ax.tick_params(axis='x', rotation=45)
    fig.tight_layout()

    canvas.draw()
    return canvas

def create_dataset(data):
    dataResponse = []
    for item in data:
        timestamp = item['timestamp']
        metric = item['metric']
        floatMetric = float(metric)
        changeRange = metric/5
        negChangeRange = changeRange * (-1)
        changeValue = metric + random.uniform(negChangeRange, changeRange)
        dataItem = {
            "timestamp":timestamp,
            "metric1":metric,
            "metric2":changeValue
                    }
        dataResponse.append(dataItem)
    print(dataResponse)
    return dataResponse
