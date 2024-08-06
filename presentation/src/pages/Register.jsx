import React, { useState, useEffect } from "react";
import { InputBox } from "../components/common/InputBox";
import { ImageSelect } from "../components/common/ImageSelect";
import { Link } from "react-router-dom";
import { SubmitButton } from "../components/common/SubmitButton";
import { SubmitForm } from "../Extra";

export const Register = ({ GlobalState }) => {
  const url = "/register";
  const [, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [passwordCheck, setPasswordCheck] = useState("");
  const [, setFirstname] = useState("");
  const [, setLastname] = useState("");
  const [, setBirth] = useState("");
  const [, setNickname] = useState("");
  const [, setAboutMe] = useState("");
  const [, setImage] = useState("");
  const { isSubmitting, setIsError } = GlobalState;

  // const [content, setContent] = useState("");
  // const [result, setResult] = useState("");
  // const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (password !== passwordCheck) {
      console.log("Passwords doesn't match!");
      setIsError("Passwords doesn't match!");
      return;
    }
    await SubmitForm(e, url, [], GlobalState); //
    // e.preventDefault();
    // setIsSubmitting(true);
    // const response = await fetch(url, {
    //   method: "POST",
    //   body: new FormData(e.target),
    // });
    // const result = await response.json();
    // setResult(result);
    // setIsSubmitting(false);
  };

  // useEffect(() => {
  //   setContent(result);
  //   if (typeof result === "object") {
  //     console.log(url, "from Server:", result);
  //   }
  // }, [result]);

  return (
    <>
      <main className="form-parent">
        <form onSubmit={handleSubmit}>
          <ImageSelect getValue={(img) => setImage(img)} />
          <InputBox
            type="text"
            id="firstname"
            label="First name"
            defaulText="First name"
            required
            getValue={(val) => setFirstname(val)}
          />
          <InputBox
            type="text"
            id="lastname"
            label="Last name"
            defaulText="Last name"
            required
            getValue={(val) => setLastname(val)}
          />
          <InputBox
            type="date"
            id="birth"
            label="Date of birth"
            defaulText=""
            getValue={(val) => setBirth(val)}
          />
          <InputBox
            type="text"
            id="nickname"
            label="Nick name"
            defaulText="nick name"
            getValue={(val) => setNickname(val)}
          />
          <InputBox
            type="email"
            id="email"
            label="Email address"
            defaulText="email address"
            required
            getValue={(val) => setEmail(val)}
          />
          <InputBox
            type="password"
            id="password"
            label="Password"
            defaulText="*********"
            required
            getValue={(val) => setPassword(val)}
          />
          <InputBox
            type="password"
            id="passwordCheck"
            label="Repeat Password"
            defaulText="*********"
            required
            getValue={(val) => setPasswordCheck(val)}
          />

          <InputBox
            type="text"
            id="aboutme"
            label="About me"
            defaulText="I am interesting because"
            getValue={(val) => setAboutMe(val)}
          />
          <SubmitButton text="Create account" disabled={isSubmitting} />
        </form>
        <div className="text-center fsize-10">
          <Link to="/login" className="btn btn-link text fsize-10">
            Allready have an account? Login in here.
          </Link>
        </div>
      </main>
    </>
  );
};
