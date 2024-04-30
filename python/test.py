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


def remove_extension(filename):
    root, _ = os.path.splitext(filename)
    return root

#TODO: Add "Match Probability" as opposed to a boolean
def process_image(image):
    passedModels = []
    for modelName in modelsArray:

        modelSavePath = f"python/savedModels/{modelName}"
        model = tf.keras.models.load_model(modelSavePath)

        # Perform inference using the loaded model
        prediction = model.predict(image)
        print("prediction for: ", model, prediction)
        predicted_label = tf.argmax(prediction, axis=1).numpy()[0]  # Assuming batch size is 1

        # Example processing: Check if prediction is a match (1) or not a match (0)
        if predicted_label == 1:
            strippedModel = remove_extension(modelName)
            passedModels.append(strippedModel)

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
    passedModels = process_image(image)  # Ensure modelsArray is defined

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
