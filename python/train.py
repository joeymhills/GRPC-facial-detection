import tensorflow as tf
from tensorflow.keras.layers import Dense, Flatten
from tensorflow.keras.models import Model
from keras.utils import to_categorical
import os
import cv2
import numpy as np
import sys

# Check if command-line arguments are provided
if len(sys.argv) < 2:
    print("Usage: python train.py personsName\n")
    sys.exit(1)  # Exit the program with an error code

modelSavePath = f"python/savedModels/{sys.argv[1]}.keras"

if os.path.exists(modelSavePath):
    sys.stderr.write("Name already taken! Please choose a new one")
    sys.exit(1)  # Exit the program with an error code

# Define num_classes
num_classes = 2

# Define paths to your image directories
train_dir = 'python/trainimg'
val_dir = 'python/testimg'
def load_images(directory):
    images = []
    labels = []
    for label in os.listdir(directory):
        label_path = os.path.join(directory, label)
        class_label = int(label.split('_')[1])  # Extract class label from directory name
        for filename in os.listdir(label_path):
            img_path = os.path.join(label_path, filename)
            img = cv2.imread(img_path)
            img = cv2.resize(img, (224, 224))  # Resize images to match input_shape
            images.append(img)
            labels.append(class_label)  # Append the correct class label, not directory name
    return np.array(images), np.array(labels)

# Load and preprocess training and validation images
x_train, y_train = load_images(train_dir)
x_val, y_val = load_images(val_dir)

# Normalize pixel values to the range [0, 1]
x_train = x_train.astype('float32') / 255.0
x_val = x_val.astype('float32') / 255.0

# Convert labels to categorical format
y_train = to_categorical(y_train, num_classes=num_classes)
y_val = to_categorical(y_val, num_classes=num_classes)

# Load a pre-trained CNN model (e.g., VGG16, ResNet50, etc.) without the top (classification) layers
base_model = tf.keras.applications.VGG16(weights='imagenet', include_top=False, input_shape=(224, 224, 3))

# Freeze the pre-trained layers so they are not updated during training
for layer in base_model.layers:
    layer.trainable = False

# Add custom layers for face recognition on top of the pre-trained model
x = Flatten()(base_model.output)
x = Dense(128, activation='relu')(x)  # Add more layers as needed

predictions = Dense(num_classes, activation='softmax')(x)  # num_classes is the number of faces you want to recognize

# Create the final model
model = Model(inputs=base_model.input, outputs=predictions)

# Compile the model
model.compile(optimizer='adam', loss='categorical_crossentropy', metrics=['accuracy'])

# Train the model using your face data
model.fit(x_train, y_train, epochs=10, batch_size=32, validation_data=(x_val, y_val))

model.save(modelSavePath)
