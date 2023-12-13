"use client";
//import Link from "next/link";
import { useEffect, useState } from "react";
import Cookies from "js-cookie";
import PrintNewMessage from "../page.jsx";


async function GetTinyProfileInfo() {
  const json = await fetch("http://localhost:8080/api/getHeadbar", {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify("voff"),
  });
  const response = await json.json();
  return response;
}

function GetTinyProfile(props) {
  const username = props.username;
  const avatar = props.avatar;
  const [msg, setMsg] = useState([]);

  return (
    <>
      <a className="fit" href="/profile">
        <div className="tinyavatar">
          <img className="pfp" src={avatar} alt="Your Avatar" />
          <div id="moo" className="link-up">
            {username}
          </div>
        </div>
      </a>
    </>
  );
}

async function getNotificationInfo() {
  const json = await fetch("http://localhost:8080/api/getNotifications", {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
  });
  const response = await json.json();
  return response;
}

function Logout() {
  Cookies.set("session_token", "value", { expires: 0, path: "/" });
}

function GetNotification(props) {
  //const [notification, setNotification] = useState(false);
  //const [number, setNumber] = useState([0]);
  const notification = props.notifications;
  const notifNumb = props.notifNumber;
  //const number = props.notifNumber
  const [number, setNumber] = useState(props.notifNumber);

  useEffect(() => {
    setNumber(notifNumb);
    //console.log("useffect entered")
  }, [notifNumb]);

  return (
    <>
      <a href="/profile" className="notification">
        <p className="notifText">{number} NOTIFICATIONS</p>
      </a>
    </>
  );
}

const Headers = (props) => {
  const newNotif = props.notifs;
  //console.log("newNotif headbar", newNotif)
  const [extraNotifs, setExtraNotifs] = useState([]);
  const [notifications, setNotifications] = useState();
  const [tinyProfInfo, setTinyProfInfo] = useState();
  const [notifNumber, setNotifNumber] = useState(0);
  const [username, setUsername] = useState([]);
  const [avatar, setAvatar] = useState([]);

  useEffect(() => {
    (async () => {
      const userdat = await GetTinyProfileInfo();
      if (userdat.status !== "success") {
        console.log("sumtin went wrong", notifdat);
      }
      setUsername(userdat.username);
      setAvatar(userdat.avatar);

      const notifdat = await getNotificationInfo();
      if (notifdat.status !== "success") {
        console.log("getting notification error", notifdat);
      }
      // console.log(notifdat.notifications, "ökasdjköas");
      if (notifdat.notifications !== null) {
        setNotifications(notifdat.notifications);
        setNotifNumber(notifdat.notifications.length);
      }
    })();
  }, []);

  useEffect(() => {
    setExtraNotifs(newNotif);
    let notifN = 0;
    let newNotifN = 0;
    if (notifications) notifN = notifications.length;
    if (newNotif) newNotifN = newNotif.length;
    setNotifNumber(notifN + newNotifN);
    //console.log("current new notifs", newNotif)
    //console.log("current new notif number", notifNumber, notifN, newNotifN)
  }, [newNotif]);

  return (
    <header className="headbar">
      <div className="identity">
        <a href="http://localhost:3000">
          <img
            className="logo"
            src="http://localhost:8080/images/Rickrolling.png"
            alt="Your Company"
          />
        </a>
      </div>
      <div>
        <GetTinyProfile username={username} avatar={avatar} />
        <a href="/group" className="fit">
          <div className="groupeButton">Group</div>
        </a>
        <a href="/chat" className="fit">
          <div className="chatButton">Chat</div>
        </a>
      </div>
      <GetNotification
        notifNumber={notifNumber}
        notifications={notifications}
      />
      <a
        href="/login"
        className="fit"
        onClick={() => {
          Logout();
        }}
      >
        <div className="logout">
          <div className="theSoundOfTaDaronne">log out</div>
        </div>
      </a>
    </header>
  );
};

export default Headers;
