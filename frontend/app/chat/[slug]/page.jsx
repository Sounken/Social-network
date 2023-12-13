"use client";
import { useEffect, useState } from "react";
import Headers from "../../components/Header";

export default function currentChat(slug) {
  const user = decodeURIComponent(slug.params.slug);
  const [chat, setChat] = useState();
  const [websocket, setWebSocket] = useState();
  const [notif, setNotif] = useState([]);
  useEffect(() => {
    (async () => {
      const c = await getChat(slug);
      setChat(c.messages);

      if (c.status === "success") {
        const newWS = new WebSocket("ws://localhost:8080/api/ws");
        newWS.onmessage = (msg) => {
          let newMsg = JSON.parse(msg.data);

          if (
            newMsg.type === "private_message" &&
            (newMsg.receiver === user || newMsg.sender === user)
          ) {
            setChat((prev) => {
              if (prev === null) {
                return [newMsg];
              }
              return [...prev, newMsg];
            });
          }
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
        setWebSocket(newWS);
        return () => {
          // console.log("closing websocket");
          newWS.close();
        };
      }
    })();
  }, []);

  return (
    <>
      <Headers notifs={notif} />
      <h2>Chatting with {user}</h2>
      <RenderChatBox message={chat} slug={slug} websocket={websocket} />
    </>
  );
}

function RenderChatBox(props) {
  const message = props.message;
  const websocket = props.websocket;
  //console.log("ws", websocket);
  const otherUser = decodeURIComponent(props.slug.params.slug);
  //console.log("render chat box", message, otherUser)
  const [messages, setMessages] = useState([]);
  const [error, setError] = useState(null);
  function handleSubmit(event) {
    event.preventDefault();

    const msg = event.target.chat.value;

    if (msg.split("").filter((c) => c !== " ").length === 0) {
      setError("Message cannot be empty");
      return;
    }

    setError(null);
    websocket.send(
      JSON.stringify({
        receiver: otherUser,
        content: msg,
        type: "private_message",
      })
    );
    event.target.chat.value = "";
  }
  useEffect(() => {
    setMessages(message);
  }, [message]);
  //setMessages(message)

  return (
    <>
      <div className="messageBox">
        <h1 className="title">Chat</h1>
        <div className="chat">
          <div className="autoScroll">
            {messages ? (
              messages.map((messages, index) => (
                <div className="messages" key={index}>
                  <p className="lineBreak">{messages.sender}</p>
                  <p className="lineBreak">{messages.content}</p>
                  <p className="lineBreak">{messages.created}</p>
                </div>
              ))
            ) : (
              <h2>No Messages</h2>
            )}
          </div>
        </div>
        <div className="chatBox">
          <form onSubmit={handleSubmit}>
            <textarea
              type="text"
              id="chat"
              name="chat"
              className="postContentCreation"
              placeholder="message"
            />
            <button type="submit">Send</button>
            {error && <p className="error">{error}</p>}
          </form>
        </div>
      </div>
    </>
  );
}

async function getChat(slug) {
  if (slug === undefined) return;

  const otherUser = decodeURIComponent(slug.params.slug);
  const json = await fetch("http://localhost:8080/api/getPrivateMessages", {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(otherUser),
  });
  const response = await json.json();
  if (response.status === "success") {
    // console.log("response", response);
    return response;
  }
  console.log("error msg", response.status);
  return response;
}

function SendMessage(reciver) {
  return (
    <>
      <textarea
        className="postContentCreation"
        type="text"
        placeholder="message"
      />
      <button type="submit">Send</button>
    </>
  );
}

export function PrintNewMessage(msg) {
  return (
    <>
      <h1>{msg.sender}</h1>
      <h1>{msg.content}</h1>
    </>
  );
}

function GetMessages(slug) {
  const user = decodeURIComponent(slug.params.slug);
  const [messages, setMessages] = useState([]);
  useEffect(() => {
    fetch("http://localhost:8080/api/getPrivateMessages", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(user),
    })
      .then((data) => data.json())
      .then((data) => {
        if (data.status !== "success") {
          console.log("error getting messages", data.status);
          return;
        }
        console.log(data.messages);
        setMessages(data.messages);
      });
  }, []);
  console.log(messages, "work");
  return (
    <>
      <div className="chatpage">
        {messages && (
          <div className="chatpage__messages">
            {messages.map((message, index) => (
              <div key={index} className="chatpage__message">
                <p>folow</p>
                <h1>{message.content}</h1>
              </div>
            ))}
          </div>
        )}
        {SendMessage(user)}
      </div>
    </>
  );
}
