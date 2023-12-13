"use client";

import { useState, useEffect, use } from "react";
import Headers from "./components/Header";
import encodeImageFile from "./components/encodeImage";
import GetYourImages from "./components/getyourimages";
import ImageSelector from "./components/imageSelector";

const datapromise = GetYourImages();

function HomePage() {
  //const img = await GetYourImages()
  const img = use(datapromise);
  // console.log("images from top", img);
  const [notif, setNotif] = useState([]);

  useEffect(() => {
    const newWS = new WebSocket("ws://localhost:8080/api/ws");
    newWS.onmessage = (msg) => {
      let newMsg = JSON.parse(msg.data);
      // console.log("new message", newMsg);
      if (
        newMsg.type === "group_join_request" ||
        newMsg.type === "group_invite" ||
        newMsg.type === "follow_request" ||
        newMsg.type === "event"
      ) {
        // console.log("new notification", newMsg);
        setNotif((prevValue) => {
          //console.log("prev",prevValue)
          //console.log("upd", [...prevValue, newMsg])

          return [...prevValue, newMsg];
        });
        //console.log("upd value", notif)
      }
    };
    //setWebSocket(newWS);
    return () => {
      // console.log("closing websocket");
      newWS.close();
    };
  }, []);
  return (
    <>
      <Headers notifs={notif} />
      <MakePost userImages={img} />
      <GetPosts />
      <GetAllUsers />
    </>
  );
}

function GetAllUsers() {
  const [users, setUsers] = useState([]);
  useEffect(() => {
    (async () => {
      fetch("http://localhost:8080/loginapi/getUsernames", {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
      })
        .then((data) => data.json())
        .then((data) => {
          if (data.status !== "success") {
            // console.log("failed to get users");
            return;
          }
          setUsers(data.users);
        });
    })();
  }, []);
  return (
    <>
      {users && (
        <ul className="userList">
          {users.map((user, index) => (
            <li key={index}>
              <a href={`profile/${user.username}`}>
                <div className="text">{user.username}</div>
              </a>
            </li>
          ))}
        </ul>
      )}
    </>
  );
}

function MakeComment(post_id) {
  const [error, setError] = useState(null);
  const [allImages, setAllImages] = useState([]);
  const [image, setImage] = useState("");
  const [content, setContent] = useState("");
  const [showImages, setShowImages] = useState(false);
  post_id = post_id.post_id;
  const post_type = "post";

  useEffect(() => {
    (async () => {
      const images = await GetYourImages();
      setAllImages(images);
    })();
  }, []);

  const handleCommentSubmit = async (e) => {
    e.preventDefault();

    if (content === "") {
      setError("Say something nice :)");
      return;
    }
    setError(null);

    fetch("http://localhost:8080/api/addComment", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ post_id, content, image, post_type }),
    })
      .then((data) => data.json())
      .then((data) => {
        // console.log("added comment", data);
        if (data.status !== "success") {
          // console.log("failed to add comment", data);
          return;
        }
        location.reload("/");
      });
  };
  return (
    <>
      <div className="makeComment">
        {image === "" ? (
          <p className="warning">No image selected</p>
        ) : (
          <img src={image} className="avatar_preview"></img>
        )}
        <form onSubmit={handleCommentSubmit}>
          <input
            type="text"
            placeholder="..."
            className="commentContentCreation"
            maxLength={255}
            onChange={(e) => setContent(e.target.value)}
          />
          {showImages ? (
            <>
              <button onClick={() => setShowImages(false)}>X</button>
              <ImageSelector images={allImages} func={setImage} />
            </>
          ) : (
            <button className="" onClick={() => setShowImages(true)}>
              Select Image
            </button>
          )}
          <button type="submit" className="commentCreationButton">
            Send
          </button>
        </form>
        <form onSubmit={encodeImageFile}>
          <input type="file" id="image" name="image"></input>
          <button type="submit" className="text">
            Upload an image
          </button>
          {error && <p className="error">{error}</p>}
        </form>
      </div>
    </>
  );
}

function ToggleComments(post_id) {
  post_id = post_id.post_id;
  const type = "normal_comments";
  const [showMore, setShowMore] = useState(false);
  const [comments, setComments] = useState([]);
  function handleClick() {
    fetch("http://localhost:8080/api/getComments", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ post_id, type }),
    })
      .then((data) => data.json())
      .then((data) => {
        if (data.status !== "success") {
          // console.log("failed to get comments", data);
        }

        setComments(data.comments);
        setShowMore(!showMore);
        // console.log(data.comments, "comments");
      });
  }
  return (
    <>
      <button className="buttonComment" onClick={handleClick}>
        comments
      </button>
      {showMore && comments && (
        <div className="commentos">
          {comments.map((dat, index) => (
            <div key={index} className="commenting">
              <div className="commentDate">{dat.created}</div>
              <a href={`profile/${dat.author}`}>
                <div className="commentUser">{dat.author}</div>
              </a>
              {dat.image_path && (
                <img src={dat.image_path} alt="image" className="pfp"></img>
              )}
              <div className="commentContent">{dat.content}</div>
            </div>
          ))}
          <MakeComment post_id={post_id} />
        </div>
      )}
      {showMore && !comments && (
        <div>
          no comments yet
          <MakeComment post_id={post_id} />
        </div>
      )}
    </>
  );
}

function MakePost({ userImages }) {
  const [error, setError] = useState(null);
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [privacy, setPrivacy] = useState("public");
  const [allowed_users, setAllowed] = useState([]);
  const [allImage, setAllImage] = useState(userImages);
  const [users, setUsers] = useState([]);
  const [image, setimage] = useState();
  const type = "post";
  const [showImages, setShowImages] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (title === "" || content === "") {
      setError("Title and Content need fulfilling");
      return;
    }
    setError(null);
    // console.log(allowed_users, "fucked");
    // console.log(image, { type, privacy, allowed_users, image, content, title });
    const res = await fetch("http://localhost:8080/api/addPost", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        type,
        privacy,
        allowed_users,
        image,
        content,
        title,
      }),
    });
    const data = await res.json();
    if (data.status !== "success") {
      // console.log("failed to make post");
      return;
    }
    // console.log("success");
    location.reload("/");
  };

  const handlePrivacyChange = async (e) => {
    setPrivacy(e);
    if (e === "semi-private") {
      const res = await fetch("http://localhost:8080/api/getUsernames", {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
      });
      const data = await res.json();
      if (data.status !== "success") {
        // console.log("failed to get users");
        return;
      }
      setUsers(data.users);
    }
  };
  const AddToAllowed = (e) => {
    if (e.target.checked) {
      setAllowed((allowed_users) => [...allowed_users, e.target.value]);
    } else {
      setAllowed((oldArray) =>
        oldArray.filter((allowed_users) => allowed_users !== e.target.value)
      );
    }
  };

  return (
    <>
      <div className="makePost">
        {image === "" ? (
          <h3>no selected image</h3>
        ) : (
          <img src={image} className="avatar_preview"></img>
        )}
        <form onSubmit={handleSubmit}>
          <input
            className="titleCreation"
            type="text"
            placeholder="title"
            maxLength="100"
            onChange={(e) => setTitle(e.target.value)}
          />
          <select
            className="dropdown"
            placeholder="private"
            onChange={(e) => handlePrivacyChange(e.target.value)}
          >
            {/*your job*/}
            <option value="public">Public</option>
            <option value="semi-private">Semi-Private</option>
            <option value="private">Private</option>
          </select>
          {privacy === "semi-private" && (
            <ul>
              {users.map((user, index) => (
                <li key={index}>
                  <input
                    type="checkbox"
                    key={index}
                    value={user.username}
                    onChange={(e) => AddToAllowed(e)}
                  />
                  <label>{user.username}</label>
                </li>
              ))}
            </ul>
          )}
          <textarea
            type="text"
            placeholder="content"
            className="postContentCreation"
            maxLength={200}
            onChange={(e) => setContent(e.target.value)}
          />
          <br />
          {error && <p className="error">{error}</p>}
          <button type="submit" className="postCreationButton">
            submit
          </button>
        </form>
        <div className="imagePoster">
          <form onSubmit={encodeImageFile}>
            <input type="file" id="image" name="image"></input>
            <button type="submit" className="text">
              AddImage
            </button>
          </form>
          {showImages ? (
            <>
              <button onClick={() => setShowImages(false)}>X</button>
              <ImageSelector images={allImage} func={setimage} />
            </>
          ) : (
            <button className="" onClick={() => setShowImages(true)}>
              Select Image
            </button>
          )}
        </div>
      </div>
    </>
  );
}

function GetPosts() {
  const [allPosts, setAllPosts] = useState([]);
  const [selectedPosts, setSelectedPosts] = useState([]);
  const [publicPost, setPublicPost] = useState([]);
  const [semiPost, setSemiPost] = useState([]);
  const [privatePost, setPrivatePost] = useState([]);

  useEffect(() => {
    fetch("http://localhost:8080/api/getPosts", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
    })
      .then((data) => data.json())
      .then((data) => {
        if (data.status !== "success") {
          // console.log("failed to get posts");
          return;
        }
        let mergedPosts = [];
        if (data.posts !== null) {
          mergedPosts = [...mergedPosts, ...data.posts];
        }
        if (data.semi_private_posts !== null) {
          mergedPosts = [...mergedPosts, ...data.semi_private_posts];
        }
        if (data.private_posts !== null) {
          mergedPosts = [...mergedPosts, ...data.private_posts];
        }
        let publicPost = data.posts;
        let semiPost = data.semi_private_posts;
        let privatePost = data.private_posts;
        if (mergedPosts !== null) {
          mergedPosts.sort(
            (a, b) => new Date(b.creation_date) - new Date(a.creation_date)
          );
          publicPost.sort(
            (a, b) => new Date(b.creation_date) - new Date(a.creation_date)
          );
          if (semiPost !== null) {
            semiPost.sort(
              (a, b) => new Date(b.creation_date) - new Date(a.creation_date)
            );
          }
          if (privatePost !== null) {
            privatePost.sort(
              (a, b) => new Date(b.creation_date) - new Date(a.creation_date)
            );
          }
        }
        setAllPosts(mergedPosts);
        setPublicPost(publicPost);
        setSemiPost(semiPost);
        setPrivatePost(privatePost);
        setSelectedPosts(mergedPosts);
      });
  }, []);

  return (
    <>
      <div>
        <button
          onClick={() => {
            setSelectedPosts(allPosts);
          }}
        >
          All
        </button>
        <button onClick={() => setSelectedPosts(publicPost)}>Public</button>
        <button onClick={() => setSelectedPosts(semiPost)}>Semi-Private</button>
        <button onClick={() => setSelectedPosts(privatePost)}>Private</button>
      </div>
      <div className="posts">
        {selectedPosts.map((post) => (
          <Post key={post.post_id} post={post}></Post>
        ))}
      </div>
    </>
  );
}

function Post({ post }) {
  return (
    <div className="post">
      <div className="postDate">
        {post.privacy} | {post.creation_date}
      </div>
      <a className="link-up" href={`profile/${post.author}`}>
        <div className="postUser">{post.author}</div>
      </a>
      <div className="postTitle">{post.title}</div>
      <div className="postContent">{post.content}</div>
      {post.image !== null &&
        post.image !== "http://localhost:8080/images/default.jpeg" &&
        post.image !== "" && (
          <img src={post.image} alt="image" className="postImage" />
        )}
      <ToggleComments post_id={post.post_id}></ToggleComments>
    </div>
  );
}

export default HomePage;
