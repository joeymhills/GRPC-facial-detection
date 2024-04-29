import tensorflow as tf
from tensorflow.keras.layers import Dense, Flatten
from tensorflow.keras.models import Model
from keras.utils import to_categorical
import os
import cv2
import numpy as np
import sys
import socket

HOST = '127.0.0.1'  # Host IP where your Go code will connect
PORT = 49522        # Port to listen on (same as your Go code)

# Check if command-line arguments are provided
if len(sys.argv) < 2:
    print("Usage: python train.py personsname\n")
    sys.exit(1)  # Exit the program with an error code

#modelsavepath = f"python/savedmodels/{sys.argv[1]}.keras"
modelSavePath = f"savedModels/lebron.keras"
#filePath = "img/temp.jpg"

filePath = "python/testimg/class_1/images.jpeg"

# Define a function to preprocess the image
def preprocess_image(received_bytes):
    received_array = np.frombuffer(received_bytes, dtype=np.uint8)

    # Decode the image array using cv2.imdecode
    img = cv2.imdecode(received_array, cv2.IMREAD_COLOR)
    img = cv2.resize(img, (224, 224))  # Resize to match the input shape of your model
    img = img.astype('float32') / 255.0  # Normalize pixel values
    return img.reshape(1, 224, 224, 3)  # Add batch dimension

def process_image(image):
    model = tf.keras.models.load_model(modelSavePath)
    # Perform inference using the loaded model
    prediction = model.predict(image)
    predicted_label = tf.argmax(prediction, axis=1).numpy()[0]  # Assuming batch size is 1

    # Example processing: Check if prediction is a match (1) or not a match (0)
    if predicted_label == 1:
        return True
    else:
        return False

#This code opens a TCP socket and waits for an image to come through
with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
    s.bind((HOST, PORT))
    s.listen()
    print("waiting for a connection")
    conn, addr = s.accept()
    with conn:
        print('Connected by', addr)
        while True:
            img_data = b''
            data = conn.recv(9000000)
            img_data += data

            # Process the received image data
            image = preprocess_image(img_data)
            # Perform image processing and get boolean response
            is_match = process_image(image)
            # Send boolean response back to Go code
            conn.sendall(is_match.to_bytes(1, 'big'))
            break
