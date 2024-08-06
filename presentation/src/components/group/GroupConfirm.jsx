import React from "react";
export const GroupConfirm = ({ id }) => {
  const handleClick = (e) => {
    console.log(id, e.target);
  };
  return (
    <div>
      <button
        title="Confirm"
        className="bi bi-person-fill-check btn btn-toggle d-inline-flex align-items-center rounded border-0 collapsed"
        onClick={handleClick}
      ></button>
    </div>
  );
};
