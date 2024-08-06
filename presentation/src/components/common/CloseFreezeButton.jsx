import React from "react";
export const CloseFreezeButton = (props) => {
  return (
    <button
      title="Close"
      className="bi bi-close btn btn-toggle d-inline-flex align-items-center rounded border-0 collapsed absolute-top-right"
      onClick={(e) => {
        props.getValue(false);
      }}
    ></button>
  );
};

// import React from "react";
// export const ConfirmButton = ({ id }) => {
// };
