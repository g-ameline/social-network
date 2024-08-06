import { useEffect, useState } from "react";
import { InputBox } from "../components/common/InputBox";
import { Link } from "react-router-dom";
import { SubmitButton } from "../components/common/SubmitButton";
import { SubmitForm } from "../Extra";

export const Login = ({ GlobalState }) => {
  const url = "/login";
  const [, setEmail] = useState("");
  const [, setPassword] = useState("");
  // const [content, setContent] = useState("");
  // const [result, setResult] = useState("");
  const { isSubmitting, setIsSubmitting, setSessionUUID, setIsError } =
    GlobalState;

  const handleSubmit = async (e) => {
    // e.preventDefault();
    await SubmitForm(e, url, [], GlobalState); //
    //   setIsSubmitting(true);
    //   const response = await fetch(url, {
    //     method: "POST",
    //     body: new FormData(e.target),
    //   });
    //   const result = await response.json();
    //   setResult(result);
    //   setIsSubmitting(false);
    //   return;
  };

  // useEffect(() => {
  //   setContent(result);
  //   if (typeof result === "object") {
  //     console.log(url, "from Server:", result);
  //     if (result.what === "login" && result.info === "Succeeded") {
  //       console.log("Should change parent getValue");
  //       getValue(true);
  //     }
  //   }
  // }, [result]);

  return (
    <>
      <main className="form-parent">
        <form onSubmit={handleSubmit}>
          <InputBox
            type="email"
            id="email"
            label="Email address"
            defaulText="name@example.com"
            required
            getValue={(val) => setEmail(val)}
          />
          <InputBox
            type="password"
            id="password"
            label="Password"
            defaulText="**********"
            required
            getValue={(val) => setPassword(val)}
          />
          <SubmitButton text="Sign in" disabled={isSubmitting} />
        </form>
        <div className="text-center fsize-10">
          <Link to="/register" className="btn btn-link text fsize-10">
            Not a member? Register
          </Link>
        </div>
      </main>
    </>
  );
};
