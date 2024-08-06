import React from "react";

export const SubmitButton = (props) => {
  return (
    <button
      className="btn btn-primary w-100 py-1"
      type="submit"
      disabled={props.disabled}
    >
      {props.text}
    </button>
  );
};
