import React, { useState } from "react";
import { GroupList } from "./GroupList";
import { GroupAdd } from "./GroupAdd";
export const GroupMenu = ({ GlobalState }) => {
  const GroupMenuList = [
    { title: "My groups", url: "/group_my", func: ["leave", "invite"] },
    {
      title: "My pending groups",
      url: "/group_pending",
      func: ["leave", "confirm"],
    },
    { title: "Other groups", url: "/group_list", func: ["join"] },
  ];
  const [isActive, setIsActive] = useState(-1);

  return (
    <div className="menu-group pt-2">
      <div className="relative">
        <h3>Group menu</h3>
        <GroupAdd title="add new group" GlobalState={GlobalState} />
      </div>
      <ul className="list-unstyled mb-0">
        {GroupMenuList.map((x, i) => {
          return (
            <li
              className={
                isActive === i
                  ? "list-group-item group-active"
                  : "list-group-item"
              }
              key={i}
              onClick={(e) => {
                e.stopPropagation();
                i === isActive ? setIsActive(-1) : setIsActive(i);
              }}
            >
              <div className="listHeader">{x.title}</div>
              {isActive === i && (
                <div className="listContent">
                  <GroupList
                    url={x.url}
                    func={x.func}
                    GlobalState={GlobalState}
                  />
                </div>
              )}
            </li>
          );
        })}
      </ul>
    </div>
  );
};
