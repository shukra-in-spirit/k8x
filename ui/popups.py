import tkinter as tk
from tkinter import ttk
import graphs
import matplotlib.pyplot as plt
from matplotlib.backends.backend_tkagg import FigureCanvasTkAgg
import requests
import time
import customtkinter
import os
from PIL import Image
import local

class ConnectPopup():
    def __init__(self):
        super().__init__()

    def open_popup_connect(self, app):

        self.main_root = app
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

        button1 = ttk.Button(popup, text="Connect", command=self.connect_to_k8s)
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

    def connect_to_k8s(self):
        self.main_root.set_status("Online - You are connected.")


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

class AddPopup():
    def __init__(self):
        super().__init__()

    def open_popup_add(self, app):
        # Create a new Toplevel window for the pop-up
        self.popup = tk.Toplevel(app)
        self.popup.title("Add New Service")

        self.main_root = app

        # Create and place form widgets in the pop-up window
        label1 = ttk.Label(self.popup, text="Deployment Name:")
        self.entry1 = ttk.Entry(self.popup)

        label2 = ttk.Label(self.popup, text="Namespace:")
        self.entry2 = ttk.Entry(self.popup)

        label3 = ttk.Label(self.popup, text="epochs:")
        entry3 = ttk.Entry(self.popup)

        label4 = ttk.Label(self.popup, text="hidden layers:")
        entry4 = ttk.Entry(self.popup)

        label5 = ttk.Label(self.popup, text="n_steps:")
        entry5 = ttk.Entry(self.popup)

        label6 = ttk.Label(self.popup, text="n_features:")
        entry6 = ttk.Entry(self.popup)


        button1 = ttk.Button(self.popup, text="Create", command=self.add_service)
        button2 = ttk.Button(self.popup, text="Predict", command=self.start_service)
        button3 = ttk.Button(self.popup, text="Cancel", command=self.popup.destroy)

        # Grid layout for form widgets
        label1.grid(row=0, column=0, padx=10, pady=5, sticky="w")
        self.entry1.grid(row=0, column=1, padx=10, pady=5, sticky="w")

        label2.grid(row=1, column=0, padx=10, pady=5, sticky="w")
        self.entry2.grid(row=1, column=1, padx=10, pady=5, sticky="w")

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

    def add_service(self):
        deployment = self.entry1.get()
        namespace = self.entry2.get()
        # url = "http://localhost:8585/" + deployment + namespace + "123"
        # requests.post(url)
        cpu_request, mem_request, container = local.get_data(deployment,namespace)
        current_dir = os.path.dirname(os.path.abspath(__file__))
        self.main_root.call_add_item(item=deployment, namespace=namespace, cpuRequest=cpu_request,memRequest=mem_request,container=container, image=customtkinter.CTkImage(Image.open(os.path.join(current_dir, "test_images", "chat_light.png"))))
        return


    def start_service(self):
        deployment = self.entry1.get()
        namespace = self.entry2.get()
        url = "http://localhost:8585/" + deployment + namespace + "123/start"
        requests.post(url)
        self.popup.destroy()
        return

class SimulationPopup():
    def __init__(self):
        super().__init__()

    def submit_form(self):
        completeDataset = graphs.create_dataset(self.dataset)
        drawnCanvas = graphs.plot_metrics_in_tkinter(completeDataset,self.canvas,self.fig,self.ax)
        return drawnCanvas

    def open_popup_simulation(self, app, item):
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


        # Plotting the data
        self.fig, self.ax = plt.subplots()
        # Embed the plot in the Tkinter window
        self.canvas = FigureCanvasTkAgg(self.fig, master=popup)
        self.dataset = [
        {"timestamp": "2023-01-01T12:00:00", "metric": 12},
        {"timestamp": "2023-01-01T13:00:00", "metric": 18},
        {"timestamp": "2023-01-01T14:00:00", "metric": 15},
        {"timestamp": "2023-01-01T15:00:00", "metric": 15},
        {"timestamp": "2023-01-01T16:00:00", "metric": 10},
        {"timestamp": "2023-01-01T17:00:00", "metric": 15},
        {"timestamp": "2023-01-01T18:00:00", "metric": 10},
        {"timestamp": "2023-01-01T19:00:00", "metric": 15},
        {"timestamp": "2023-01-01T20:00:00", "metric": 10},
        {"timestamp": "2023-01-01T21:00:00", "metric": 15},
        {"timestamp": "2023-01-01T22:00:00", "metric": 10},
        {"timestamp": "2023-01-01T23:00:00", "metric": 15},
        {"timestamp": "2023-01-02T00:00:00", "metric": 10},
        {"timestamp": "2023-01-02T01:00:00", "metric": 15},
        {"timestamp": "2023-01-02T02:00:00", "metric": 10},
        {"timestamp": "2023-01-02T03:00:00", "metric": 15},
        # Add more data as needed
    ]

        button2 = ttk.Button(popup, text="Simulate", command=self.submit_form)
        button3 = ttk.Button(popup, text="Cancel", command=popup.destroy)

        label3.grid(row=0, column=0, padx=10, pady=5, sticky="w")
        entry3.grid(row=0, column=1, padx=10, pady=5, sticky="w")

        label4.grid(row=1, column=0, padx=10, pady=5, sticky="w")
        entry4.grid(row=1, column=1, padx=10, pady=5, sticky="w")

        label5.grid(row=2, column=0, padx=10, pady=5, sticky="w")
        entry5.grid(row=2, column=1, padx=10, pady=5, sticky="w")

        label6.grid(row=3, column=0, padx=10, pady=5, sticky="w")
        entry6.grid(row=3, column=1, padx=10, pady=5, sticky="w")
        self.canvas.get_tk_widget().grid(row=4, column=0, columnspan=2)

        button2.grid(row=5, column=0, columnspan=2, pady=5)
        button3.grid(row=6, column=0, columnspan=2, pady=5)