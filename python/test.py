import tensorflow as tf
import os
import cv2
import numpy as np
import sys
import socket
import threading

# Check if command-line arguments are provided
if len(sys.argv) < 2:
    print("Usage: python train.py personsname\n")
    sys.exit(1)  # Exit the program with an error code


def process_arguments(arg_string):
    # Split the argument string by commas and strip whitespace
    arguments = [arg.strip() for arg in arg_string.split(',')]
    return arguments

modelsArray = process_arguments(sys.argv[1])


# Define a function to preprocess the image
def preprocess_image(received_bytes):
    received_array = np.frombuffer(received_bytes, dtype=np.uint8)

    # Decode the image array using cv2.imdecode
    img = cv2.imdecode(received_array, cv2.IMREAD_COLOR)
    img = cv2.resize(img, (224, 224))
    img = img.astype('float32') / 255.0  # Normalize pixel values
    return img.reshape(1, 224, 224, 3)  # Add batch dimension

def remove_extension(model_name):
    return model_name.replace(".keras", "")

# Runs image against each of the CNNs
def process_image(image):
    passedModels = []

    # List to store (modelName, prediction) tuples
    model_predictions = []

    for modelName in modelsArray:
        modelSavePath = f"python/savedModels/{modelName}"
        model = tf.keras.models.load_model(modelSavePath)

        # Perform inference using the loaded model
        prediction = model.predict(image)
        print("prediction for:", modelName, prediction)
        predicted_label = tf.argmax(prediction, axis=1).numpy()[0]

        # Append (modelName, prediction) tuple to the list
        model_predictions.append((modelName, prediction[0, 1]))  # Assuming prediction is a 2D array

    # Sort the list of tuples based on the second element (prediction value)
    sorted_models = sorted(model_predictions, key=lambda x: x[1], reverse=True)

    # Print sorted predictions for debugging
    print("Sorted predictions:", sorted_models)

    # Extract the model names from the sorted list of tuples
    passedModels = [model_tuple[0] for model_tuple in sorted_models if model_tuple[1] == 1.0]
    
    # Remove ".keras" suffix from model names
    passedModels = [remove_extension(model_name) for model_name in passedModels]

    return passedModels


def handle_client(conn, addr):
    print('Connected by', addr)
    while True:
        img_data = b''
        data = conn.recv(9000000)
        img_data += data
        break

    # Process the received image data
    image = preprocess_image(img_data)

    # Perform image processing and get boolean response
    passedModels = process_image(image)

    # Convert list of passed models to a comma-separated string
    response = ','.join(passedModels)
    conn.sendall(response.encode())  # Send response back to Go code
    conn.close()


def start_server():
    host = '127.0.0.1'
    port = 49522
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.bind((host, port))
        s.listen()
        print("\n\nWaiting for connections...")
        while True:
            conn, addr = s.accept()
            threading.Thread(target=handle_client, args=(conn, addr)).start()


start_server()
