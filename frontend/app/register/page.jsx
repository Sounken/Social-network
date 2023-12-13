"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import cookie from "js-cookie";

export default function Register() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [birthDate, setBirthDate] = useState("");
  const [error, setError] = useState("");

  const router = useRouter();

  const handleSubmit = async (e) => {
    e.preventDefault();
    const response = await fetch("http://localhost:8080/api/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email, password, firstName, lastName, birthDate }),
    });
    const data = await response.json();

    //setRegistrationDone(true) and hide the mandatory form, show the optional form
    if (data.status === "success") {
      cookie.set("session_token", data.token);
      router.push("/optional");
    } else {
      setError("The email is already in use");
      return;
    }
  };

  return (
    <>
      <header className="headbar">
        <img
          className="logo"
          src="http://localhost:8080/images/Rickrolling.png"
          alt="Your Company"
        />
        <h1 className="padder">Irrelevent Discussion</h1>
      </header>
      <div className="signin-window">
        <form className="organ" onSubmit={handleSubmit}>
          {/*Email*/}
          <div className="gibspace">
            <label htmlFor="email">Email address</label>
            <div className="mt-2">
              <input
                id="email"
                name="email"
                type="email"
                placeholder="Email@Example"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="input-field"
              ></input>
            </div>
          </div>
          {/*First Name*/}
          <div className="gibspace">
            <label htmlFor="firstName">First Name</label>
            <div className="mt-2">
              <input
                id="firstName"
                name="firstName"
                type="firstName"
                placeholder="First Name"
                required
                value={firstName}
                onChange={(e) => setFirstName(e.target.value)}
                className="input-field"
              />
            </div>
          </div>
          {/* Last Name */}
          <div className="gibspace">
            <label htmlFor="lastName">Last Name</label>
            <div className="mt-2">
              <input
                id="lastName"
                name="lastName"
                type="lastName"
                placeholder="Last Name"
                required
                value={lastName}
                onChange={(e) => setLastName(e.target.value)}
                className="input-field"
              />
            </div>
          </div>
          {/* Birth Date */}
          <div className="gibspace">
            <label htmlFor="birthDate">Date of Birth</label>
            <div className="mt-2">
              <input
                id="birthDate"
                name="birthDate"
                type="date"
                min="1899-01-01"
                max="2023-01-01"
                required
                value={birthDate}
                onChange={(e) => setBirthDate(e.target.value)}
                className="input-field"
              />
            </div>
          </div>
          {/* Password */}
          <div className="gibspace">
            <label htmlFor="password">Password</label>
            <div className="mt-2">
              <input
                id="password"
                name="password"
                type="password"
                placeholder="Password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="input-field"
              />
            </div>
          </div>
          {/* Submit Button */}
          <div className="gibspace">
            {error && <p className="error">{error}</p>}
            <button type="submit" className="padder">
              Register
            </button>
          </div>
        </form>
        <p className="text">
          Already a user?{" "}
          <Link href="/login" className="link-up">
            Login here
          </Link>
        </p>
      </div>
    </>
  );
}
