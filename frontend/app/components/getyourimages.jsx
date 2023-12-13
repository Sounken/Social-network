async function GetYourImages() {
  const resp = await fetch("http://localhost:8080/api/getYourImages", {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
  });
  const data = await resp.json();
  if (data.status !== "success") {
    console.log("failed to get your images", data.status);
    return;
  }
  //console.log("got your images", data)
  return data;
}

export default GetYourImages;
