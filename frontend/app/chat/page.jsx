"use client";
import React, { useState, useEffect, use } from "react";
import Headers from "../components/Header";

function AllChats() {
  const [users, setUsers] = useState([]);
  useEffect(() => {
    (async () => {
      fetch("http://localhost:8080/api/getUsernames", {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
      })
        .then((res) => res.json())
        .then((data) => {
          if (data.status !== "success") {
          }
          setUsers(data.users);
        });
    })();
  }, []);

  return (
    <>
      <div className="container">
        {users.length > 0 && (
          <div className="chatorganise">
            <h2 className="title">Chat</h2>
            {users.map((user, index) => (
              <div className="user" key={index}>
                <h1>
                  <a className="contents" href={`/chat/${user.username}`}>
                    {user.username}
                  </a>
                </h1>
              </div>
            ))}
          </div>
        )}
      </div>
    </>
  );
}

function Chatpage() {
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
        if (preValue !== null) {
          setNotif((prevValue) => [...prevValue, newMsg]);
        }
      }
    };

    return () => {
      newWS.close();
    };
  }, []);
  return (
    <>
      <Headers notifs={notif} />
      <AllChats />
    </>
  );
}

export default Chatpage;
