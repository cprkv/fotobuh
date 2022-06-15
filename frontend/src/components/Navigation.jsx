import React from "react";
import { NavLink } from "react-router-dom";
import styled from "styled-components";
import { apiGetCategories } from "../api";

const NavigationList = styled.nav`
  padding: 1rem 0.3rem 1.3rem 1rem;
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-start;
  align-items: center;
`;

const NavigationLink = styled(NavLink)`
  padding: 0.3rem 1rem 0.5rem 0;
  align-items: center;
  display: flex;

  &.active {
    font-weight: bold;
  }
`;

const NavigationIcon = styled.img`
  width: 1.3rem;
  margin: 0.09rem 0.3rem 0;
`;

export function Navigation() {
  const [categories, setCategories] = React.useState([]);
  React.useEffect(() => {
    apiGetCategories().then(setCategories);
  }, []);

  return (
    <NavigationList>
      <NavigationLink to="/" key="index">
        <NavigationIcon src="/imgs/icon.svg" />
        latest
      </NavigationLink>

      {categories && categories.map((cat) => (
        <NavigationLink
          to={`/category/${cat.id}`}
          key={cat.id}
          activeclassname="active"
        >
          {cat.name}
        </NavigationLink>
      ))}
    </NavigationList>
  );
}
