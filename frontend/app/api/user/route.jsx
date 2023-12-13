import { NextResponse } from "next/server";
import user from "./user.json";

export async function GET(request) {
  return NextResponse.json(user);
}
