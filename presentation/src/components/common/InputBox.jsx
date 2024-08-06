import React, { useState } from "react";

export const InputBox = (props) => {
  const [elem, setElem] = useState("");
  return (
    <div className="form-floating">
      <input
        type={props.type}
        className="form-control height-25rem mb-3"
        id={props.id}
        name={props.id}
        placeholder={props.defaulText}
        required={"required" in props && "required"}
        value={elem}
        onChange={(e) => {
          setElem(e.target.value);
          props.getValue(e.target.value);
        }}
      />
      <label className="fsize-10" htmlFor={props.id}>
        {props.label}
      </label>
    </div>
  );
};
