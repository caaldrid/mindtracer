import React from "react";
import { NavLink, Outlet } from "react-router-dom";

const Layout: React.FC = () => {
  const navItems = [
    { label: "Projects", path: "/projects" },
    { label: "Areas", path: "/areas" },
    { label: "Resources", path: "/resources" },
    { label: "Archive", path: "/archive" },
  ];

  return (
    <div>
      <nav>
        {navItems.map((item) => (
          <NavLink
            key={item.path}
            to={item.path}
          >
            {item.label}
          </NavLink>
        ))}
      </nav>

      <main>
        <Outlet />
      </main>
    </div>
  );
};

export default Layout;
