function ImageSelector(props) {
  const images = props.images;
  const updateSelectedImage = props.func;

  return (
    <>
      <div>
        <button onClick={() => updateSelectedImage("")}>No Image</button>
      </div>
      <div className="shroud">
        {images.stock_images ? (
          images.stock_images.map((image, index) => (
            <div className="text" key={index}>
              <img
                src={image}
                alt="image"
                className="pfp"
                onClick={() => updateSelectedImage(image)}
              />
            </div>
          ))
        ) : (
          <h1>no stock images</h1>
        )}
        <br />
        {images.user_images ? (
          images.user_images.map((image, index) => (
            <div className="text" key={index}>
              <img
                src={image}
                alt="image"
                className="pfp"
                onClick={() => updateSelectedImage(image)}
              />
            </div>
          ))
        ) : (
          <h1></h1>
        )}
      </div>
    </>
  );
}

export default ImageSelector;
