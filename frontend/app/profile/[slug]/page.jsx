"use client";

import { useState, useEffect } from "react";
import Headers from "../../components/Header";
import { useRouter } from "next/navigation";

function FollowCheck({ slug }) {
  const [following, setFollwing] = useState(Boolean);
  const user = decodeURIComponent(slug.params.slug);
  useEffect(() => {
    fetch("http://localhost:8080/api/followCheck", {
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
          // console.log("follow check failed");
          //return;
        }
        setFollwing(data.following);
      });
  }, []);

  const handleFollow = () => {
    fetch("http://localhost:8080/api/addOrRemoveFollow", {
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
          // console.log(data.status);
          // console.log("unfollow/follow failed");
          return;
        }
        // console.log("unfollow/follow success 123213 12", user);
      });
  };
  return (
    <>
      {following ? (
        <>
          <a href={`${user}`}>
            <button onClick={() => handleFollow()}>Unfollow</button>
          </a>
        </>
      ) : (
        <>
          <a href={`${user}`}>
            <button onClick={() => handleFollow()}>Follow</button>
          </a>
        </>
      )}
    </>
  );
}

function GetProfile({ slug }) {
  const user = decodeURIComponent(slug.params.slug);
  const [stuff, setStuff] = useState([]);
  const [followers, setFollowers] = useState([]);
  const [following, setFollowing] = useState([]);
  const [groups, setGroups] = useState([]);
  const router = useRouter();
  const [isPrivate, setIsPrivate] = useState(true);

  useEffect(() => {
    fetch("http://localhost:8080/api/getProfile", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(user),
    })
      .then((data) => data.json())
      .then((data) => {
        // console.log(data, "data");
        if (data.status !== "success") {
          // console.log("get profile failed", data); // Update the state with fetched data
        }
        if (data.you === true) {
          router.push("/profile");

          //return;
        }
        if (data.status !== "private") {
          setIsPrivate(false);
        }
        setStuff(data);
        setFollowers(data.followers);
        setFollowing(data.following);
        setGroups(data.groups);
      });
  }, []);
  return (
    <>
      {!isPrivate || following ? (
        <>
          <div className="Profile">
            <div className="organiser">
              <img className="avatar_preview" src={stuff.avatar} />
            </div>
            <div className="docu">
              <p> {stuff.username}</p>
              <p> {stuff.first_name}</p>
              <p> {stuff.last_name}</p>
              <p> Email: {stuff.email}</p>
              <p> Bio: {stuff.bio}</p>
              <p> Brithday: {stuff.dob}</p>
              <p> {stuff.privacy}</p>
              <FollowCheck slug={slug} />
            </div>
          </div>
          <div>
            <div className="folow">
              <h2>Followers</h2>
              {followers && (
                <div>
                  {followers.map((follower, index) => (
                    <a key={index} href={`/profile/${follower}`}>
                      <p>{follower}</p>
                    </a>
                  ))}
                </div>
              )}

              {following && (
                <div>
                  <h2>Following</h2>
                  {following.map((follow, index) => (
                    <div key={index}>
                      <a
                        className="link-up"
                        key={index}
                        href={`/profile/${follow}`}
                      >
                        <p>{follow}</p>{" "}
                      </a>
                    </div>
                  ))}
                </div>
              )}

              {groups && (
                <div>
                  <h2>Groups</h2>
                  {groups.map((group, index) => (
                    <div key={index}>
                      <p>{group}</p>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </>
      ) : (
        <>
          <div>
            <h2>Profile</h2>
            <div>{stuff.email}</div>
            <FollowCheck slug={slug} />
          </div>
        </>
      )}
    </>
  );
}

function ProfilePage(slug) {
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
      console.log("closing websocket");
      newWS.close();
    };
  }, []);
  return (
    <>
      <Headers notifs={notif} />
      <div className="layouter">
        <GetProfile slug={slug} />
      </div>
    </>
  );
}

export default ProfilePage;
