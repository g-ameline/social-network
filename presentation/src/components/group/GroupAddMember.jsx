import React, { useState } from "react";
import { SearchBox } from "../common/SearchBox";
import { SubmitButton } from "../common/SubmitButton";
import { SubmitForm } from "../../Extra";
export const GroupAddMember = ({ id, GlobalState }) => {
  const [member, setMember] = useState("");
  const handleSubmit = (e) => {
    e.preventDefault();
    const url = "/group_request";
    const formData = new FormData();
    formData.append("group_id", id);
    formData.append("session", GlobalState.SessionUUID);
    formData.append("user", member);
    console.log("Andmed:", e, url, formData, GlobalState);
    SubmitForm(e, url, formData, GlobalState);
  };
  return (
    <form onSubmit={handleSubmit}>
      <SearchBox
        GlobalState={GlobalState}
        getValue={(val) => {
          setMember(val);
        }}
      />
      <SubmitButton text="Add member" />
    </form>
  );
};
