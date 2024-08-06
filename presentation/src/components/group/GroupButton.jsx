import React from "react";
import { SubmitForm } from "../../Extra";
export const GroupButton = ({
  id,
  GlobalState,
  url,
  title,
  clss,
  xtra,
  value,
  setValue,
}) => {
  const handleClick = async (e) => {
    if (typeof setValue !== "undefined") {
      value === true ? setValue(false) : setValue(true);
      return;
    }
    const formData = new FormData();
    formData.append("group_id", id);
    formData.append("session", GlobalState.SessionUUID);
    if (xtra !== "") {
      formData.append("who", xtra);
    }
    SubmitForm(e, url, formData, GlobalState);
  };
  return (
    <div>
      <button
        title={title}
        className={
          "bi btn btn-toggle d-inline-flex align-items-center rounded border-0 collapsed " +
          clss
        }
        onClick={handleClick}
      ></button>
    </div>
  );
};
