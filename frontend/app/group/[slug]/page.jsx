"use client";

import Headers from "../../components/Header";
import { useRouter } from "next/navigation";
import GetYourImages from "../../components/getyourimages";
import encodeImageFile from "../../components/encodeImage";
import { useState, useEffect, use } from "react";
import ImageSelector from "@/app/components/imageSelector";
import Link from "next/link";

function InviteToGroup(props) {
  const slug = props.slug;
  const members = props.members;
  // console.log("members invToGroup", members);
  const group_name = decodeURIComponent(slug.params.slug);
  const [username, setUsername] = useState("");
  const [allusers, setAllUsers] = useState([]);
  const [invite_success, setInviteSuccess] = useState(false);

  const filterOutMembers = (user) => {
    for (let i = 0; i < members.length; i++) {
      if (
        user.email === members[i].username ||
        user.username === members[i].username
      ) {
        return false;
      }
    }
    return true;
  };

  useEffect(() => {
    (async () => {
      fetch("http://localhost:8080/api/getUsernames", {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
      })
        .then((data) => data.json())
        .then((data) => {
          if (data.status !== "success") {
            console.log("failed to get usernames", data.status);
            return;
          }

          setAllUsers(data.users.filter(filterOutMembers));
        });
    })();
  }, []);
  // console.log("all users", allusers);
  const handleChange = (e) => {
    setUsername(e.target.value);
    setInviteSuccess(false);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const receiver = username;
    fetch("http://localhost:8080/api/sendGroupInvite", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ receiver, group_name }),
    })
      .then((data) => data.json())
      .then((data) => {
        if (data.status !== "success") {
          console.log("failed to invite user", data.status);
          return;
        }
        setInviteSuccess(true);
        // console.log("invited user", data.status);
      });
  };
  if (invite_success) {
    return (
      <div>
      <form onSubmit={handleSubmit}>
        {allusers && (
          <select onChange={(e) => handleChange(e)}>
            {allusers.map((user, index) => (
              <option key={index} value={user.username}>
                {user.username}
              </option>
            ))}
          </select>
        )}
        <p>{username} has been invited</p>
        <button type="submit">YES</button>
      </form>
    </div>
    );
  }
  return (
    <>
      <div>
        <form onSubmit={handleSubmit}>
          {allusers && (
            <select onChange={(e) => setUsername(e.target.value)}>
              {allusers.map((user, index) => (
                <option key={index} value={user.username}>
                  {user.username}
                </option>
              ))}
            </select>
          )}
          <p> Do you want to invite {username}</p>
          <button type="submit">YES</button>
        </form>
      </div>
    </>
  );
}

async function addMemberToGroupChat(slug) {
  if (slug === undefined) return "not_yet";
  const groupName = decodeURIComponent(slug.params.slug);

  const json = await fetch("http://localhost:8080/api/addGroupMemberToChat", {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(groupName),
  });
  const data = await json.json();
  return data.status;
}

async function getGroupPageData(slug) {
  if (slug === undefined) return;
  const groupname = decodeURIComponent(slug.params.slug);
  const json = await fetch("http://localhost:8080/api/getGroupPage", {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(groupname),
  });
  const data = await json.json();
  if (data.status !== "success") {
    // console.log("error in get group", data.status);
    return data.status;
  }

  return data;
}

//const promisdata2 = getGroupPageData()

//const promisdata = GroupPage()
export default function GroupPage(slug) {
  const [isMember, setIsMember] = useState(false);
  const [chat, setChat] = useState([]);
  const [groupPageData, setGroupPageData] = useState("not_a_member");
  const [groupExists, setGroupExists] = useState(true);
  const [websocket, setWebSocket] = useState(null);
  const [notif, setNotif] = useState([]);

  useEffect(() => {
    (async () => {
      const m = await addMemberToGroupChat(slug);
      if (m !== "success") {
        if (m === "group_does_not_exist") {
          setGroupExists(false);
          return;
        }
      }
      const c = await ChatBox(slug);
      setChat(c);
      const g = await getGroupPageData(slug);
      setGroupPageData(g);
      if (g !== "not_a_member") {
        setIsMember(true);
        const newWS = new WebSocket("ws://localhost:8080/api/ws");
        newWS.onmessage = (msg) => {
          //console.log("new notification",msg)
          let newMsg = JSON.parse(msg.data);
          if (newMsg.type === "group_message") {
            // console.log("new notification parsed", newMsg);

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
      if (g === "group_does_not_exist") {
        setGroupExists(false);
      }
    })();
  }, []);
  return (
    <>
      {isMember && groupExists ? (
        <>
          <Headers notifs={notif} />
          <div className="standards">
            <InviteToGroup slug={slug} members={groupPageData.members} />
            <MakeGroupPost slug={slug} />
          </div>
          <RenderGroup data={groupPageData} slug={slug} />
          <RenderChatBox message={chat} slug={slug} websocket={websocket} />
        </>
      ) : groupExists ? (
        <>
          <Headers notifs={notif} />
          <RenderGroup data={groupPageData} slug={slug} />
        </>
      ) : (
        <>
          <Headers notifs={notif} />
          <h1>Group does not exist</h1>
        </>
      )}
    </>
  );
}

function RenderGroup(props) {
  const data = props.data;
  const slug = props.slug;
  const [groupPosts, setGroupPosts] = useState([]);
  const [events, setEvents] = useState([]);
  const [members, setMembers] = useState([]);
  const [joinRequest, setJoinRequest] = useState([]);
  const [userType, setUserType] = useState("");
  const [isMember, setIsMember] = useState(false);
  const groupname = decodeURIComponent(slug.params.slug);

  useEffect(() => {
    if (data !== "not_a_member") {
      setGroupPosts(data.group_posts);
      setEvents(data.events);
      setMembers(data.members);
      setJoinRequest(data.join_request);
      setUserType(data.member_type);
      setIsMember(true);
    } else {
      setUserType("not_a_member");
      setIsMember(false);
    }
  }, []);

  // console.log(events, "events");
  return (
    <>
      <div>
        {isMember ? (
          <div className="groupPage">
            <h1>{groupname}</h1>

            {groupPosts && (
              <div className="container">
                {groupPosts.map((post, index) => (
                  <div className="groupPost" key={index}>
                    <div className="postti">
                      <div className="poster">
                        <h2>{post.author}</h2>
                      </div>

                      {post.image !== null &&
                        post.image !==
                          "http://localhost:8080/images/default.jpeg" &&
                        post.image !== "" && (
                          <img
                            src={post.image}
                            alt="image"
                            className="postImage"
                          />
                        )}

                      <div className="contents">
                        <p>{post.content}</p>
                        <p>{post.creation_date}</p>
                      </div>
                      {/* <MakeComment post_id={post.post_id}></MakeComment> */}
                      <ToggleComments
                        post_id={post.post_id}
                        group_name={groupname}
                      ></ToggleComments>
                    </div>
                  </div>
                ))}
              </div>
            )}
            {members.username && (
              <div>
                {members.username && members.username.length == 1 && (
                  <div>
                    <h1>Members</h1>
                    <p>{members.username}</p>
                  </div>
                )}
                {members.map((member, index) => (
                  <div className="groupMember" key={index}>
                    <div className="groupMemberInfo">
                      <h2>{member.username}</h2>
                      <h1> {member.member_type}</h1>
                    </div>

                    <div className="groupMemberInfo">
                      <h2>{member.email}</h2>
                      <h1> {member.member_type}</h1>
                    </div>
                  </div>
                ))}
              </div>
            )}
            <h2>Event page</h2>
            <MakeEvent slug={slug} />
            {/* MakeEvent(slug) */}
            {events && (
              <div className="container">
                {events.map((event, index) => (
                  <div className="groupEvent" key={index}>
                    <Event event={event} groupName={groupname} />
                    
                  </div>
                ))}
              </div>
            )}
            {joinRequest && (
              <div className="container">
                {joinRequest.map((request, index) => (
                  <div className="groupPost" key={index}>
                    {request && (
                      <div className="postti">
                        <h2>{request.username}</h2>
                        <h1>{request.email}</h1>
                      </div>
                    )}
                  </div>
                ))}
              </div>
            )}

            {userType === "creator" && (
              <div className="container">
                <h1>Join Requests</h1>
                {joinRequest && (
                  <div>
                    {joinRequest.map((request, index) => (
                      <div className="groupJoinRequest" key={index}>
                        {request && (
                          <div className="groupJoinRequestInfo">
                            <h2>{request.username}</h2>
                            <h1>{request.email}</h1>
                          </div>
                        )}
                      </div>
                    ))}
                  </div>
                )}
              </div>
            )}
          </div>
        ) : (
          <>
            <h2>{groupname}</h2>
            <button onClick={() => RequestToJoin(slug)}>Request to join</button>
          </>
        )}
      </div>
    </>
  );
}

function RenderChatBox(props) {
  const message = props.message;
  const websocket = props.websocket;
  const groupName = decodeURIComponent(props.slug.params.slug);
  const [messages, setMessages] = useState([]);
  const [error, setError] = useState();

  function handleSubmit(event) {
    event.preventDefault();

    if (event.target.chat.value === "") {
      setError("Please enter a message");
      return;
    }
    setError(null);

    const msg = event.target.chat.value;
    websocket.send(
      JSON.stringify({
        receiver: groupName,
        content: msg,
        type: "group_message",
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
            {messages && (
              <div>
                {messages.map((messages, index) => (
                  <div className="messages" key={index}>
                    <p className="lineBreak">{messages.sender}</p>
                    <p className="lineBreak">{messages.content}</p>
                    <p className="lineBreak">{messages.created}</p>
                  </div>
                ))}
              </div>
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
              placeholder="Enter your message"
            />
            <button type="submit">Send</button>
            {error && <p className="error">{error}</p>}
          </form>
        </div>
      </div>
    </>
  );
}

function Event(props) {
  const event = props.event;
  const groupname = props.groupName;
  const [hasAnswered, setHasAnswered] = useState(event.already_chosen);

  const HandleEventResponse = async (answer, event_id) => {
    const group_name = groupname;
    fetch("http://localhost:8080/api/handleEventAnswer", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ group_name, answer, event_id }),
    })
      .then((data) => data.json())
      .then((data) => {
        if (data.status !== "success") {
          // console.log("failed to handle event answer", data.status);
          return;
        }
        // console.log("handled event answer", data.status);
      });
  };

  return (
    <>
      <div className="groupPost">
        <div className="postti">
          <h2>{event.title}</h2>
          <p>{event.content}</p>

          <div className="contents">
            <p>{event.date}</p>
            <p>{event.event_date}</p>
          </div>
          {!hasAnswered ? (
            <form>
              <button
                onClick={() => {
                  HandleEventResponse("going", event.event_id);
                }}
              >
                Going
              </button>
              <button
                onClick={() => {
                  HandleEventResponse("not going", event.event_id);
                }}
              >
                Not Going
              </button>
            </form>
          ) : (
            event.options.map((opt, ind) => (
              <p key={ind}>
                {opt}: {event.answers[ind]}
              </p>
            ))
          )}
        </div>
      </div>
    </>
  );
}

async function ChatBox(slug) {
  if (slug === undefined) {
    return;
  }
  const groupname = decodeURIComponent(slug.params.slug);
  const json = await fetch("http://localhost:8080/api/getGroupChat", {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(groupname),
  });
  const response = await json.json();
  if (response.status === "success") {
    return response.messages;
  }
  return;
}

function MakeGroupPost(props) {
  const slug = props.slug;
  const group_name = decodeURIComponent(slug.params.slug);
  const type = "group_post";
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [image, setImage] = useState("");
  const [allImage, setAllImage] = useState([]);
  const [showImages, setShowImages] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    (async () => {
      const images = await GetYourImages();
      setAllImage(images);
    })();
  }, []);

  const SubmitHandle = async (e) => {
    e.preventDefault();

    if (title === "" || content === "") {
      setError("Please fill in all fields");
      return;
    }
    setError(null);

    fetch("http://localhost:8080/api/addPost", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ group_name, title, content, type, image }),
    })
      .then((data) => data.json())
      .then((data) => {
        if (data.status !== "success") {
          // console.log("failed to add post", data);
          return;
        }
        // console.log("added post", data);
      });
  };
  return (
    <>
      <div>
        {image === "" ? (
          <p className="warning">no selected image</p>
        ) : (
          <img src={image} className="avatar_preview"></img>
        )}
        <form onSubmit={encodeImageFile}>
          <input type="file" id="image" name="image"></input>
          <button type="submit" className="text">
            Add Image
          </button>
        </form>
      </div>
      <div>
        <form className="groupPostmaker" onSubmit={SubmitHandle}>
          <input
            className="titleCreation"
            type="text"
            placeholder="..."
            onChange={(e) => setTitle(e.target.value)}
          />

          <textarea
            className="postContentCreation"
            type="text"
            placeholder="what do you want to say?"
            onChange={(e) => setContent(e.target.value)}
          />
          {error && <p className="error">{error}</p>}
          <button type="submit">submit</button>
        </form>
        {showImages ? (
          <>
            <button onClick={() => setShowImages(false)}>X</button>
            <ImageSelector images={allImage} func={setImage} />
          </>
        ) : (
          <button className="" onClick={() => setShowImages(true)}>
            Select Image
          </button>
        )}
      </div>
    </>
  );
}

function MakeEvent({ slug }) {
  const router = useRouter();
  const group_name = decodeURIComponent(slug.params.slug);
  const [title, setEventTitle] = useState("");
  const [content, setContent] = useState("");
  const [event_date, setEvent_Date] = useState("");
  const [error, setError] = useState(null);
  const options = ["going", "not going"]

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (title === "" || content === "" || event_date === "") {
      setError("Fill in all fields!");
      return;
    }
    setError(null);
    fetch("http://localhost:8080/api/addEvent", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ group_name, title, content, event_date, options }),
    })
      .then((data) => data.json())
      .then((data) => {
        if (data.status !== "success") {
          // console.log("failed to add event", data);
          return;
        }
        // console.log(data);
      });
    location.reload();
  };
  return (
    <>
      <div className="signin-window">
        <form onSubmit={handleSubmit}>
          <input
            type="text"
            placeholder="Event Name"
            className="mt-2"
            onChange={(e) => setEventTitle(e.target.value)}
          />
          <textarea
            type="text"
            className="postContentCreation"
            placeholder="Event Description"
            onChange={(e) => setContent(e.target.value)}
          />
          <input
            type="date"
            className="mt-2"
            placeholder="event date"
            onChange={(e) => setEvent_Date(e.target.value)}
          />

          <button type="submit">Create Event</button>
          {error && <p className="error">{error}</p>}
        </form>
      </div>
    </>
  );
}

function RequestToJoin(slug) {
  const group_name = decodeURIComponent(slug.params.slug);
  fetch("http://localhost:8080/api/sendGroupJoinRequest", {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ group_name }),
  })
    .then((data) => data.json())
    .then((data) => {
      if (data.status !== "success") {
        // console.log(data.status);
        // console.log("request to join failed");
        return;
      }
      // console.log("request to join success");
    });
}

function MakeComment(props) {
  const post_id = props.post_id;
  const group_name = props.group_name;
  const [image, setImage] = useState("");
  const [allImage, setAllImage] = useState([]);
  const [showImages, setShowImages] = useState(false);
  const [content, setContent] = useState("");
  const [error, setError] = useState(null);
  const post_type = "group_post";

  useEffect(() => {
    (async () => {
      const images = await GetYourImages();
      setAllImage(images);
    })();
  }, []);

  const handleCommentSubmit = async (e) => {
    e.preventDefault();
    if (content === "") {
      setError("Say some nice!");
      return;
    }
    setError(null);
    fetch("http://localhost:8080/api/addComment", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ post_id, content, image, post_type, group_name }),
    })
      .then((data) => data.json())
      .then((data) => {
        if (data.status !== "success") {
          // console.log("failed to add comment", data);
          return;
        }
        // console.log("added comment", data);
      });
  };
  return (
    <>
      <div className="makeComment">
        {image === "" ? (
          <p className="warning">no selected image</p>
        ) : (
          <img src={image} className="avatar_preview"></img>
        )}

        <form onSubmit={handleCommentSubmit}>
          <input
            type="text"
            placeholder="..."
            className="commentContentCreation"
            onChange={(e) => setContent(e.target.value)}
          />
          <button type="submit" className="commentCreationButton">
            Send
          </button>
          {error && <p className="error">{error}</p>}
        </form>
        {showImages ? (
          <>
            <button onClick={() => setShowImages(false)}>X</button>
            <ImageSelector images={allImage} func={setImage} />
          </>
        ) : (
          <button className="" onClick={() => setShowImages(true)}>
            Select Image
          </button>
        )}
        <form onSubmit={encodeImageFile}>
          <input type="file" id="image" name="image"></input>
          <button type="submit" className="text">
            AddImage
          </button>
        </form>
      </div>
    </>
  );
}

function ToggleComments(props) {
  const post_id = props.post_id;
  const group_name = props.group_name;
  const [showMore, setShowMore] = useState(false);
  const [comments, setComments] = useState([]);
  const type = "group_comments";
  function handleClick() {
    fetch("http://localhost:8080/api/getComments", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ post_id, type, group_name }),
    })
      .then((data) => data.json())
      .then((data) => {
        if (data.status !== "success") {
          // console.log("failed to get comments", data);
        }

        setComments(data.comments);
        setShowMore(!showMore);
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
              <Link href={`profile/${dat.author}`}>
                <div className="commentUser">{dat.author}</div>
              </Link>
              {dat.image_path && (
                <img src={dat.image_path} alt="image" className="pfp"></img>
              )}
              <div className="commentContent">{dat.content}</div>
            </div>
          ))}
          <MakeComment post_id={post_id} group_name={group_name} />
        </div>
      )}
      {showMore && !comments && (
        <div>
          no comments here yet
          <MakeComment post_id={post_id} group_name={group_name} />
        </div>
      )}
    </>
  );
}
