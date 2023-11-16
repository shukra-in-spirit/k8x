import customtkinter
import os
from PIL import Image
import popups
import time


def AuditLogs(log):
    with open("../logs.txt", 'a+') as f:
        f.write(log + "\n")

class ScrollableLabelButtonFrame(customtkinter.CTkScrollableFrame):
    def __init__(self, master, commandSimulate=None, commandStop=None, **kwargs):
        super().__init__(master, **kwargs)
        self.grid_columnconfigure(0, weight=1)
        self.label_list = []
        self.req_label_list = []      #new
        self.container_label_list = [] #new
        self.ns_label_list = []        #new
        self.button_list = []
        self.commandStop=commandStop
        self.commandSimulate=commandSimulate
        self.button_counter=0
        self.label_counter=0

    def add_item(self, item, namespace, cpuRequest, memRequest, container, image=None):
        time.sleep(2)
        label = customtkinter.CTkLabel(self, text=item, image=image, compound="left", padx=15, anchor="w", font=customtkinter.CTkFont(size=14, weight="bold"))
        namespaceString = "ns: " + namespace
        namespaceLabel = customtkinter.CTkLabel(self, text=namespaceString, padx=15, anchor="w")
        requestString = "cpu: " + cpuRequest + ", mem: " + memRequest
        requestLabel = customtkinter.CTkLabel(self, text=requestString, padx=15, anchor="w")
        containerString = "container: " + container
        containerLabel = customtkinter.CTkLabel(self, text=containerString, padx=15, anchor="w")
        buttonSimulate = customtkinter.CTkButton(self, text="Simulate", width=100, height=24)
        buttonStop = customtkinter.CTkButton(self, text="Stop", width=100, height=24)
        if self.commandStop is not None:
            buttonStop.configure(command=lambda: self.commandStop(item))
        if self.commandSimulate is not None:
            buttonSimulate.configure(command=lambda: self.commandSimulate(item))
        label.grid(row=len(self.label_list), column=0, pady=(20, 30), sticky="w")
        namespaceLabel.grid(row=len(self.label_list), column=1, pady=(20, 30), sticky="w")
        requestLabel.grid(row=len(self.label_list), column=2, pady=(20, 30), sticky="w")
        containerLabel.grid(row=len(self.label_list), column=3, pady=(20, 30), sticky="w")
        buttonSimulate.grid(row=len(self.label_list), column=4, pady=(20, 30), padx=(0,10))
        buttonStop.grid(row=len(self.label_list), column=5, pady=(20, 30), padx=(0,30))
        self.label_list.append(label)
        self.button_list.append(buttonSimulate)
        self.button_list.append(buttonStop)
        self.container_label_list.append(containerLabel)
        self.req_label_list.append(requestLabel)
        self.ns_label_list.append(namespaceLabel)
        print(self.label_list)
        print(self.button_list)

    def remove_item(self, item):
        for label in self.label_list:
            button_simulate = self.button_list[self.button_counter]
            self.button_counter=self.button_counter + 1
            button_stop = self.button_list[self.button_counter]
            self.button_counter=self.button_counter + 1
            ns_label = self.ns_label_list[self.label_counter]
            req_label = self.req_label_list[self.label_counter]
            container_label = self.container_label_list[self.label_counter]
            self.label_counter = self.label_counter + 1
            if label.cget("text") == item:
                print(self.button_counter)
                label.destroy()
                self.label_list.remove(label)
                button_simulate.destroy()
                self.button_list.remove(button_simulate)
                button_stop.destroy()
                self.button_list.remove(button_stop)
                ns_label.destroy()
                self.ns_label_list.remove(ns_label)
                req_label.destroy()
                self.req_label_list.remove(req_label)
                container_label.destroy()
                self.container_label_list.remove(container_label)
        self.button_counter=0
        return
    
    def open_simulation_panel(self, app, item):
        simulationPopup = popups.SimulationPopup()
        simulationPopup.open_popup_simulation(app, item)
        return


class App(customtkinter.CTk):
    def __init__(self):
        super().__init__()

        self.title("KubeChi - AI controlled kubernetes scaling")
        self.grid_rowconfigure(0, weight=1)
        self.columnconfigure(0, weight=0)
        self.columnconfigure(1, weight=1)

        # create sidebar frame with widgets
        self.sidebar_frame = customtkinter.CTkFrame(self, width=140, corner_radius=0)
        self.sidebar_frame.grid(row=0, column=0, rowspan=4, sticky="nsew")
        self.sidebar_frame.grid_rowconfigure(4, weight=1)
        self.logo_label = customtkinter.CTkLabel(self.sidebar_frame, text="k8x", font=customtkinter.CTkFont(size=20, weight="bold"))
        self.logo_label.grid(row=0, column=0, padx=20, pady=(20, 10))
        self.sidebar_button_1 = customtkinter.CTkButton(self.sidebar_frame, text="Connect", command=self.connect_button_event)
        self.sidebar_button_1.grid(row=1, column=0, padx=20, pady=10)
        self.sidebar_button_2 = customtkinter.CTkButton(self.sidebar_frame, text="Add Service", command=self.add_button_event)
        self.sidebar_button_2.grid(row=2, column=0, padx=20, pady=10)
        self.sidebar_button_1 = customtkinter.CTkButton(self.sidebar_frame, text="Audit Logs", command=self.log_button_event)
        self.sidebar_button_1.grid(row=3, column=0, padx=20, pady=10)
        self.sidebar_button_3 = customtkinter.CTkButton(self.sidebar_frame, text="Help", command=self.help_button_event)
        self.sidebar_button_3.grid(row=4, column=0, padx=20, pady=10)
        self.appearance_mode_label = customtkinter.CTkLabel(self.sidebar_frame, text="Appearance Mode:", anchor="w")
        self.appearance_mode_label.grid(row=5, column=0, padx=20, pady=(10, 0))
        self.appearance_mode_optionemenu = customtkinter.CTkOptionMenu(self.sidebar_frame, values=["Light", "Dark", "System"],
                                                                       command=self.change_appearance_mode_event)
        self.appearance_mode_optionemenu.grid(row=6, column=0, padx=20, pady=(10, 10))
        self.scaling_label = customtkinter.CTkLabel(self.sidebar_frame, text="UI Scaling:", anchor="w")
        self.scaling_label.grid(row=7, column=0, padx=20, pady=(10, 0))
        self.scaling_optionemenu = customtkinter.CTkOptionMenu(self.sidebar_frame, values=["80%", "90%", "100%", "110%", "120%"],
                                                               command=self.change_scaling_event)
        self.scaling_optionemenu.grid(row=8, column=0, padx=20, pady=(10, 20))

        # create scrollable label and button frame
        current_dir = os.path.dirname(os.path.abspath(__file__))
        self.status_label = customtkinter.CTkLabel(self, text="Offline - Not Connected", anchor="w")
        self.status_label.grid(row=0, column=2, padx=20, pady=(10, 0))
        self.scrollable_label_button_frame = ScrollableLabelButtonFrame(master=self, width=900, height=500,commandStop=self.label_button_frame_stop_event,commandSimulate=self.label_button_frame_simulate_event, corner_radius=0)
        self.scrollable_label_button_frame.grid(row=1, column=2, padx=10, pady=10, sticky="nsew")
        # for i in range(20):  # add items with images
        #     self.scrollable_label_button_frame.add_item(f"container number {i}", image=customtkinter.CTkImage(Image.open(os.path.join(current_dir, "test_images", "chat_light.png"))))
    def set_status(self, status):
        time.sleep(1)
        self.status_label.configure(text=status)

    def call_add_item(self, item, namespace, cpuRequest, memRequest, container, image):
        self.scrollable_label_button_frame.add_item(item, namespace, cpuRequest, memRequest, container, image)

    def label_button_frame_simulate_event(self, item):
        self.scrollable_label_button_frame.open_simulation_panel(app, item)

    def label_button_frame_stop_event(self, item):
        self.scrollable_label_button_frame.remove_item(item)

    def change_appearance_mode_event(self, new_appearance_mode: str):
        customtkinter.set_appearance_mode(new_appearance_mode)

    def change_scaling_event(self, new_scaling: str):
        new_scaling_float = int(new_scaling.replace("%", "")) / 100
        customtkinter.set_widget_scaling(new_scaling_float)

    def connect_button_event(self):
        print("sidebar_button click")
        connectPopup = popups.ConnectPopup()
        connectPopup.open_popup_connect(app)
    
    def add_button_event(self):
        print('add service button clicked')
        addPopup = popups.AddPopup()
        addPopup.open_popup_add(app)
    
    def log_button_event(self):
        print('add service button clicked')

    def help_button_event(self):
        print('add service button clicked')
        popups.open_popup_help(app)


if __name__ == "__main__":
    customtkinter.set_appearance_mode("dark")
    app = App()
    app.mainloop()