"use client";
import { NextRequest, NextResponse } from "next/server";

export async function middleware(req, NextRequest) {
  const cookiee = req.cookies.get("session_token");
  console.log("should be session", cookiee);
  let response = NextResponse.next();

  if (cookiee === undefined) {
    if (
      req.nextUrl.href === "http://localhost:3000/login" ||
      req.nextUrl.href === "http://localhost:3000/register"
    ) {
      return response;
    } else return NextResponse.redirect("http://localhost:3000/login");
  }

  const tryFetch = async (url) => {
    return fetch(url, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
        Cookie: `${cookiee.name}=${cookiee.value}`,
      },
    })
      .then((fetchResponse) => fetchResponse.json())
      .then((data) => {
        if (data.status !== "success") {
          const responso = NextResponse.redirect("http://localhost:3000/login");
          responso.cookies.delete("session_token");
          return responso;
        }
        if (
          data.status === "success" &&
          (req.nextUrl.href === "http://localhost:3000/login" ||
            req.nextUrl.href === "http://localhost:3000/register")
        ) {
          return NextResponse.redirect("http://localhost:3000/");
        }
        console.log("cookie check success");
        return response;
      });
  };

  try {
    return await tryFetch("http://localhost:8080/api/cookie");
  } catch (error) {
    console.log("Error fetching from localhost, trying backend: ", error);
    return tryFetch("http://backend:8080/api/cookie");
  }
}

export const config = {
  matcher: [
    "/",
    "/post",
    "/profile/:path*",
    "/optional",
    "/msg/:path*",
    "/group/:path*",
    "/chat/:path*",
    "/login",
    "/register",
  ],
};
