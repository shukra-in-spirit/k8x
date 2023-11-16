import tkinter as tk
from tkinter import ttk
import graphs
import matplotlib.pyplot as plt
from matplotlib.backends.backend_tkagg import FigureCanvasTkAgg

def open_popup_connect(app):
    # Create a new Toplevel window for the pop-up
    popup = tk.Toplevel(app)
    popup.title("Connect to k8x")

    # Create and place form widgets in the pop-up window
    label1 = ttk.Label(popup, text="url:")
    entry1 = ttk.Entry(popup)

    label2 = ttk.Label(popup, text="username:")
    entry2 = ttk.Entry(popup)

    label3 = ttk.Label(popup, text="password:")
    entry3 = ttk.Entry(popup)

    label4 = ttk.Label(popup, text="port:")
    entry4 = ttk.Entry(popup)

    button1 = ttk.Button(popup, text="Connect", command=submit_form)
    button2 = ttk.Button(popup, text="Cancel", command=popup.destroy)

    # Grid layout for form widgets
    label1.grid(row=0, column=0, padx=10, pady=5, sticky="w")
    entry1.grid(row=0, column=1, padx=10, pady=5, sticky="w")

    label2.grid(row=1, column=0, padx=10, pady=5, sticky="w")
    entry2.grid(row=1, column=1, padx=10, pady=5, sticky="w")

    label3.grid(row=2, column=0, padx=10, pady=5, sticky="w")
    entry3.grid(row=2, column=1, padx=10, pady=5, sticky="w")

    label4.grid(row=3, column=0, padx=10, pady=5, sticky="w")
    entry4.grid(row=3, column=1, padx=10, pady=5, sticky="w")

    button1.grid(row=4, column=0, columnspan=2, pady=10)
    button2.grid(row=5, column=0, columnspan=2, pady=5)

def submit_form():
    # Placeholder function for form submission
    print("Form submitted!")

def open_popup_help(app):
    # Create a new Toplevel window for the pop-up
    popup = tk.Toplevel(app)
    popup.title("Connect to k8x")

    # Create and place form widgets in the pop-up window
    label1 = ttk.Label(popup, text="url:")
    entry1 = ttk.Entry(popup)

    label2 = ttk.Label(popup, text="username:")
    entry2 = ttk.Entry(popup)

    label3 = ttk.Label(popup, text="password:")
    entry3 = ttk.Entry(popup)

    label4 = ttk.Label(popup, text="port:")
    entry4 = ttk.Entry(popup)

    button1 = ttk.Button(popup, text="Connect", command=submit_form)
    button2 = ttk.Button(popup, text="Cancel", command=popup.destroy)

    # Grid layout for form widgets
    label1.grid(row=0, column=0, padx=10, pady=5, sticky="w")
    entry1.grid(row=0, column=1, padx=10, pady=5, sticky="w")

    label2.grid(row=1, column=0, padx=10, pady=5, sticky="w")
    entry2.grid(row=1, column=1, padx=10, pady=5, sticky="w")

    label3.grid(row=2, column=0, padx=10, pady=5, sticky="w")
    entry3.grid(row=2, column=1, padx=10, pady=5, sticky="w")

    label4.grid(row=3, column=0, padx=10, pady=5, sticky="w")
    entry4.grid(row=3, column=1, padx=10, pady=5, sticky="w")

    button1.grid(row=4, column=0, columnspan=2, pady=10)
    button2.grid(row=5, column=0, columnspan=2, pady=5)


def open_popup_add(app):
    # Create a new Toplevel window for the pop-up
    popup = tk.Toplevel(app)
    popup.title("Add New Service")

    # Create and place form widgets in the pop-up window
    label1 = ttk.Label(popup, text="Deployment Name:")
    entry1 = ttk.Entry(popup)

    label2 = ttk.Label(popup, text="Namespace:")
    entry2 = ttk.Entry(popup)

    label3 = ttk.Label(popup, text="epochs:")
    entry3 = ttk.Entry(popup)

    label4 = ttk.Label(popup, text="hidden layers:")
    entry4 = ttk.Entry(popup)

    label5 = ttk.Label(popup, text="n_steps:")
    entry5 = ttk.Entry(popup)

    label6 = ttk.Label(popup, text="n_features:")
    entry6 = ttk.Entry(popup)


    button1 = ttk.Button(popup, text="Create", command=submit_form)
    button2 = ttk.Button(popup, text="Predict", command=submit_form)
    button3 = ttk.Button(popup, text="Cancel", command=popup.destroy)

    # Grid layout for form widgets
    label1.grid(row=0, column=0, padx=10, pady=5, sticky="w")
    entry1.grid(row=0, column=1, padx=10, pady=5, sticky="w")

    label2.grid(row=1, column=0, padx=10, pady=5, sticky="w")
    entry2.grid(row=1, column=1, padx=10, pady=5, sticky="w")

    label3.grid(row=2, column=0, padx=10, pady=5, sticky="w")
    entry3.grid(row=2, column=1, padx=10, pady=5, sticky="w")

    label4.grid(row=3, column=0, padx=10, pady=5, sticky="w")
    entry4.grid(row=3, column=1, padx=10, pady=5, sticky="w")

    label5.grid(row=4, column=0, padx=10, pady=5, sticky="w")
    entry5.grid(row=4, column=1, padx=10, pady=5, sticky="w")

    label6.grid(row=5, column=0, padx=10, pady=5, sticky="w")
    entry6.grid(row=5, column=1, padx=10, pady=5, sticky="w")

    button1.grid(row=6, column=0, columnspan=2, pady=10)
    button2.grid(row=7, column=0, columnspan=2, pady=5)
    button3.grid(row=8, column=0, columnspan=2, pady=5)

def open_popup_simulation(app, item):
    # Create a new Toplevel window for the pop-up
    popup = tk.Toplevel(app)
    popup.title(item)

    label3 = ttk.Label(popup, text="epochs:")
    entry3 = ttk.Entry(popup)

    label4 = ttk.Label(popup, text="hidden layers:")
    entry4 = ttk.Entry(popup)

    label5 = ttk.Label(popup, text="n_steps:")
    entry5 = ttk.Entry(popup)

    label6 = ttk.Label(popup, text="n_features:")
    entry6 = ttk.Entry(popup)

    button2 = ttk.Button(popup, text="Simulate", command=submit_form)
    button3 = ttk.Button(popup, text="Cancel", command=popup.destroy)

    # Plotting the data
    fig, ax = plt.subplots()
    # Embed the plot in the Tkinter window
    canvas = FigureCanvasTkAgg(fig, master=popup)
    dataset = [
    {"timestamp": "2023-01-01T12:00:00", "metric1": 10, "metric2": 20},
    {"timestamp": "2023-01-01T13:00:00", "metric1": 15, "metric2": 25},
    {"timestamp": "2023-01-01T14:00:00", "metric1": 10, "metric2": 20},
    {"timestamp": "2023-01-01T15:00:00", "metric1": 15, "metric2": 25},
    {"timestamp": "2023-01-01T16:00:00", "metric1": 10, "metric2": 20},
    {"timestamp": "2023-01-01T17:00:00", "metric1": 15, "metric2": 25},
    {"timestamp": "2023-01-01T18:00:00", "metric1": 10, "metric2": 20},
    {"timestamp": "2023-01-01T19:00:00", "metric1": 15, "metric2": 25},
    {"timestamp": "2023-01-01T20:00:00", "metric1": 10, "metric2": 20},
    {"timestamp": "2023-01-01T21:00:00", "metric1": 15, "metric2": 25},
    {"timestamp": "2023-01-01T22:00:00", "metric1": 10, "metric2": 20},
    {"timestamp": "2023-01-01T23:00:00", "metric1": 15, "metric2": 25},
    {"timestamp": "2023-01-01T00:00:00", "metric1": 10, "metric2": 20},
    {"timestamp": "2023-01-01T01:00:00", "metric1": 15, "metric2": 25},
    {"timestamp": "2023-01-01T02:00:00", "metric1": 10, "metric2": 20},
    {"timestamp": "2023-01-01T03:00:00", "metric1": 15, "metric2": 25},
    # Add more data as needed
]
    drawnCanvas = graphs.plot_metrics_in_tkinter(dataset,canvas,fig,ax)


    label3.grid(row=0, column=0, padx=10, pady=5, sticky="w")
    entry3.grid(row=0, column=1, padx=10, pady=5, sticky="w")

    label4.grid(row=1, column=0, padx=10, pady=5, sticky="w")
    entry4.grid(row=1, column=1, padx=10, pady=5, sticky="w")

    label5.grid(row=2, column=0, padx=10, pady=5, sticky="w")
    entry5.grid(row=2, column=1, padx=10, pady=5, sticky="w")

    label6.grid(row=3, column=0, padx=10, pady=5, sticky="w")
    entry6.grid(row=3, column=1, padx=10, pady=5, sticky="w")
    drawnCanvas.get_tk_widget().grid(row=4, column=0, columnspan=2)

    button2.grid(row=5, column=0, columnspan=2, pady=5)
    button3.grid(row=6, column=0, columnspan=2, pady=5)