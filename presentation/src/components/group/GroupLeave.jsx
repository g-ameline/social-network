import React from "react";
import { SubmitForm } from "../../Extra";

export const GroupLeave = ({ id, GlobalState }) => {
  const url = "/group_leave";
  const handleClick = async (e) => {
    const formData = new FormData();
    formData.append("group_id", id);
    formData.append("session", GlobalState.SessionUUID);
    SubmitForm(e, url, formData, GlobalState);
  };
  return (
    <div>
      <button
        title="Leave group"
        className="bi bi-person-fill-dash btn btn-toggle d-inline-flex align-items-center rounded border-0 collapsed"
        onClick={handleClick}
      ></button>
    </div>
  );
};
