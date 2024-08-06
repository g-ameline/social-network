import React, { useEffect, useState } from "react";

export const ImageSelect = (props) => {
  const [image, setImage] = useState("");
  const [preview, setPreview] = useState("./uploads/person.png");
  useEffect(() => {
    if (image) {
      const reader = new FileReader();
      reader.onloadend = () => {
        setPreview(reader.result);
      };
      reader.readAsDataURL(image);
    } else {
      setPreview("./uploads/person.png");
    }
  }, [image]);
  return (
    <div>
      <div className="d-flex justify-content-center">
        <div>
          <label
            title="Choose image"
            className="btn btn-link text mb-2 pt-0 position-relative"
            htmlFor="image"
          >
            <div className="d-flex justify-content-center">
              <img
                src={preview}
                className="rounded-circle"
                alt="Add icon"
                width="120"
                height="80"
              />
            </div>
            {image === "" && (
              <div className="choose-image-txt">Choose image</div>
            )}
          </label>
          <input
            type="file"
            className="form-control d-none"
            id="image"
            name="image"
            accept="image/*"
            onChange={(e) => {
              const file = e.target.files[0];
              if (file && file.type.substring(0, 5) === "image") {
                var reader = new FileReader();
                var rawData = new ArrayBuffer();
                reader.loadend = function () {};
                reader.onload = function (e) {
                  rawData = e.target.result;
                  console.log(rawData);
                };
                setImage(file);
                props.getValue(file);
              }
            }}
          />
        </div>
      </div>
    </div>
  );
};
