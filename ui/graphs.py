import tkinter as tk
from datetime import datetime
import json

def plot_metrics_in_tkinter(data, canvas, fig, ax):
    # Extract data from the JSON list
    timestamps = [entry['timestamp'] for entry in data]
    metric1_values = [entry['metric1'] for entry in data]
    metric2_values = [entry['metric2'] for entry in data]

    # Convert timestamps to datetime objects
    timestamps = [datetime.fromisoformat(timestamp) for timestamp in timestamps]

    ax.plot(timestamps, metric1_values, label='Metric 1', marker='o')
    ax.plot(timestamps, metric2_values, label='Metric 2', marker='o')

    # Configure the plot
    ax.set_xlabel('Timestamp')
    ax.set_ylabel('Values')
    ax.set_title('Comparison of Metric 1 and Metric 2 over Time')
    ax.legend()
    ax.tick_params(axis='x', rotation=45)
    fig.tight_layout()

    canvas.draw()
    return canvas
