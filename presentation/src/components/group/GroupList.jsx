import React, { useState, useEffect } from "react";
// import { GroupLeave } from "./GroupLeave";
// import { GroupInvite } from "./GroupInvite";
import { GroupButton } from "./GroupButton";
import { SubmitForm } from "../../Extra";
import { GroupAddMember } from "./GroupAddMember";

let imgUrl =
  window.location.protocol + "//" + window.location.host + "/uploads/";

export const GroupList = ({ url, func, GlobalState }) => {
  const [isFetching, setIsFetching] = useState(true);
  const [GroupData, setGroupData] = useState("");
  const [inviteToGroup, setInviteToGroup] = useState({});
  useEffect(() => {
    setIsFetching(true);
    console.log("url:", url);
    const fetchData = async (e) => {
      const formData = new FormData();
      formData.append("session", GlobalState.SessionUUID);
      setGroupData(await SubmitForm(e, url, formData, GlobalState));
      setIsFetching(false);
    };
    fetchData();
  }, []);

  if (isFetching) {
    return <div>loading...</div>;
  }
  return (
    <ul id="pending_group" className="list-group">
      {GroupData.info.map((x, i) => (
        <li
          key={i}
          prop={x.group_id}
          className="list-group-item"
          title={x.group_desc}
          onClick={(e) => {
            e.stopPropagation();
            console.log("click on sub");
          }}
        >
          <div className="d-flex">
            <div className="me-auto">
              <img
                className="list-group-image"
                src={imgUrl + x.group_icon}
                // alt={imgUrl + x.group_icon}
                alt="IMG"
              />
              <div
                className="middle"
                prop={x.group_id}
                onClick={async (e) => {
                  console.log(e.target);
                  const url = "/page_content";
                  const formData = new FormData();
                  formData.append("session", GlobalState.SessionUUID);
                  formData.append("group_id", e.target.val);
                  await SubmitForm(e, url, formData, GlobalState);
                }}
              >
                {x.group_name}
              </div>
            </div>
            {func.includes("confirm") && (
              <GroupButton
                id={x.group_id}
                GlobalState={GlobalState}
                url="/group_approve"
                title="Confirm group"
                clss="bi-person-fill-check"
                xtra="user"
              />
            )}
            {func.includes("invite") && (
              <GroupButton
                id={x.group_id}
                title="Add new members"
                clss={
                  inviteToGroup[x.group_id] === true
                    ? "bi-close-black"
                    : "bi-person-fill-add"
                }
                value={inviteToGroup[x.group_id]}
                setValue={(val) => {
                  const inviteGrp = { ...inviteToGroup, [x.group_id]: val };
                  setInviteToGroup(inviteGrp);
                }}
                GlobalState={GlobalState}
                xtra="other"
              />
            )}
            {func.includes("join") && (
              <GroupButton
                id={x.group_id}
                GlobalState={GlobalState}
                url="/group_request"
                title="Join group"
                clss="bi-person-fill-add"
                xtra="user"
              />
            )}
            {func.includes("leave") && (
              <GroupButton
                id={x.group_id}
                GlobalState={GlobalState}
                url="/group_leave"
                title="Leave group"
                clss="bi-person-fill-dash"
              />
            )}
          </div>
          {inviteToGroup[x.group_id] === true && (
            <GroupAddMember id={x.group_id} GlobalState={GlobalState} />
          )}
        </li>
      ))}
    </ul>
  );
};
