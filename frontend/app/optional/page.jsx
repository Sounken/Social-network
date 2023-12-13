"use client";

import { useState, useEffect, use } from "react";
import { useRouter } from "next/navigation";
import Headers from "../components/Header";
import encodeImageFile from "../components/encodeImage";
import GetYourImages from "../components/getyourimages";

import Link from "next/link";

function Avatars({ arg }) {
  const [stockImages, setStockImages] = useState([]);
  const [userImages, setUserImages] = useState([]);
  const images = GetYourImages();
  useEffect(() => {
    (async () => {
      const images = await GetYourImages();
      setUserImages(images.user_images);
      setStockImages(images.stock_images);
    })();
  }, []);

  return (
    <>
      <div className="padder">
        {stockImages.map((image, index) => (
          <img
            src={image}
            onClick={(e) => arg(image)}
            key={index}
            alt={`Avatar ${index}`}
            className="pfp"
          />
        ))}
      </div>
      {userImages && (
        <div className="padder">
          {userImages.map((image, index) => (
            <img
              src={image}
              onClick={(e) => arg(image)}
              key={index}
              alt={`Avatar ${index}`}
              className="pfp"
            />
          ))}
        </div>
      )}
    </>
  );
}

export default function Optional() {
  const [nickname, setNickname] = useState("");
  const [aboutMe, setAboutMe] = useState("");
  const [privacy, setPrivate] = useState("");
  const [avatar, setAvatar] = useState("");
  const router = useRouter();
  const [notif, setNotif] = useState([]);

  useEffect(() => {
    const newWS = new WebSocket("ws://localhost:8080/api/ws");
    newWS.onmessage = (msg) => {
      let newMsg = JSON.parse(msg.data);
      console.log("new message", newMsg);
      if (
        newMsg.type === "group_join_request" ||
        newMsg.type === "group_invite" ||
        newMsg.type === "follow_request" ||
        newMsg.type === "event"
      ) {
        console.log("new notification", newMsg);
        setNotif((prevValue) => [...prevValue, newMsg]);
      }
    };
    //setWebSocket(newWS);
    return () => {
      // console.log("closing websocket");
      newWS.close();
    };
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    const response = await fetch("http://localhost:8080/api/updateSettings", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ nickname, aboutMe, avatar, privacy }),
    });
    const data = await response.json();
    //setRegistrationDone(true) and hide the mandatory form, show the optional form
    if (data.status === "success") {
      router.push("/");
    }
  };
  return (
    <>
      <Headers notifs={notif} />
      <div className="signin-window">
        <form onSubmit={encodeImageFile}>
          <input type="file" id="image" name="image"></input>
          <button type="submit" className="text">
            Upload an image
          </button>
        </form>
        <form className="sorting" onSubmit={handleSubmit}>
          <div className="padder">
            <img src={avatar} alt="" value={0} className="avatar_preview" />
          </div>
          <div className="container">
            <a type="submit" id={0} className="pick-me-profile">
              <Avatars arg={setAvatar}></Avatars>
            </a>
          </div>
          <div>
            <label id="image">Choose an image:</label>
          </div>
          <br />
          <label id="aboutMe">About Me</label>
          <div>
            <textarea
              id="aboutMe"
              name="aboutMe"
              type="aboutMe"
              placeholder="About Me"
              value={aboutMe}
              onChange={(e) => setAboutMe(e.target.value)}
              className="aboutMe"
            ></textarea>
          </div>
          <br />
          <label id="nickname">Nickname</label>
          <div className="mt-2">
            <input
              id="nickname"
              name="nickname"
              type="nickname"
              placeholder="nickname"
              value={nickname}
              onChange={(e) => setNickname(e.target.value)}
              className="titleCreation"
            ></input>
          </div>
          <div>
            <input
              type="radio"
              id="public"
              name="access"
              value="public"
              onClick={(e) => setPrivate(e.target.value)}
            ></input>
            <label id="public">Public</label>
            <br></br>
            <input
              type="radio"
              id="private"
              name="access"
              value="private"
              onClick={(e) => setPrivate(e.target.value)}
            ></input>
            <label id="private">Private</label>
            <br></br>
          </div>

          <div className="finish">
            <button type="submit" className="postCreationButton">
              Finish
            </button>
          </div>
        </form>

        <p className="text">
          don't wanna expose yourself?{" "}
          <Link href="/" className="link-up">
            Skip
          </Link>
        </p>
      </div>
    </>
  );
}
