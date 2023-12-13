"use client";

function encodeImageFile(element) {
  let file = element.target.image.files[0];

  if (file === undefined) return;
  let reader = new FileReader();
  reader.readAsDataURL(file);
  reader.onloadend = function () {
    fetch("http://localhost:8080/api/addImage", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(reader.result),
    });
  };
}

export default encodeImageFile;
