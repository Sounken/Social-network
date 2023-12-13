"use client";
import React, { useEffect, useState } from "react";
import Headers from "../components/Header";
import { useRouter } from "next/navigation";

function MakeGroup() {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const router = useRouter();
  const [error, setError] = useState(null);

  const handleNewGroup = (e) => {
    e.preventDefault();

    if (name === "") {
      setError("Name cannot be empty");
      return;
    }
    setError(null);

    fetch("http://localhost:8080/api/addGroup", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ name, description }),
    })
      .then((data) => data.json())
      .then((data) => {
        if (data.status !== "success") {
          // console.log("failed to make group");
          return;
        }
        // console.log("group made");
        const url = `/group/${name}`;
        router.push(url);
      });
  };
  return (
    <>
      <form className="groupmaker" onSubmit={(e) => handleNewGroup(e)}>
        <input
          className="titleCreation"
          type="text"
          placeholder="Group Name"
          onChange={(e) => setName(e.target.value)}
        ></input>
        <br />
        <textarea
          className="aboutMe"
          type="text"
          placeholder="Group Description"
          onChange={(e) => setDescription(e.target.value)}
        ></textarea>
        <br />
        <button type="submit" className="postCreationButton">
          Make Group
        </button>
        {error && <p className="error">{error}</p>}
      </form>
    </>
  );
}

function YourGroups() {
  const [groups, setGroups] = useState([]);
  useEffect(() => {
    fetch("http://localhost:8080/api/getGroupnames", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
    })
      .then((data) => data.json())
      .then((data) => {
        if (data.status !== "success") {
          // console.log("error");
          return;
        }
        setGroups(data.groups);
      });
  }, []);
  return (
    <>
      <div className="yourGroups">
        <h2>Your Groups</h2>
        {groups.map((group, index) => (
          <div className="group" key={index}>
            <a className="link-up" href={`group/${group}`}>
              {group}
            </a>
          </div>
        ))}
      </div>
    </>
  );
}

function Groups() {
  const [websocket, setWebSocket] = useState(null);
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
        setNotif((prevValue) => [...prevValue, newMsg]);
      }
    };
    //setWebSocket(newWS);
    return () => {
      // console.log("closing websocket");
      newWS.close();
    };
  }, []);
  //console.log("notif", notif)

  return (
    <>
      <Headers notifs={notif} />
      <MakeGroup />
      <YourGroups />
    </>
  );
}

export default Groups;
