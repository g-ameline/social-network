import React from "react";
export const GroupInvite = ({ value, setValue }) => {
  return (
    <div>
      <button
        title="Invite to group"
        className={
          value === true
            ? "bi bi-close-black btn btn-toggle d-inline-flex align-items-center rounded border-0 collapsed"
            : "bi bi-person-fill-add btn btn-toggle d-inline-flex align-items-center rounded border-0 collapsed"
        }
        onClick={() => {
          value === true ? setValue(false) : setValue(true);
        }}
      ></button>
    </div>
  );
};
