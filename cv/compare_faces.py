import face_recognition
import sys

def compare_faces(image_path1, image_path2):
    # Load the first image
    image1 = face_recognition.load_image_file(image_path1)
    # Load the second image
    image2 = face_recognition.load_image_file(image_path2)
    
    # Get the face encodings for the faces in each image
    face_encodings1 = face_recognition.face_encodings(image1)
    face_encodings2 = face_recognition.face_encodings(image2)
    
    # If no faces are found in either image, return False
    if len(face_encodings1) == 0 or len(face_encodings2) == 0:
        return False
    
    # Compare the faces found in the first image to the faces found in the second image
    for face_encoding1 in face_encodings1:
        matches = face_recognition.compare_faces(face_encodings2, face_encoding1, tolerance=0.6)
        if True in matches:
            return True
    
    return False

if __name__ == "__main__":
    image_path1 = sys.argv[1]
    image_path2 = sys.argv[2]

    result = compare_faces(image_path1, image_path2)
    if result:
        print("true")
    else:
        print("false")