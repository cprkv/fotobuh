import React from "react";
import styled from "styled-components";

const FooterWrapper = styled.footer`
  padding: 1rem;
  text-align: center;
`;

export function Footer() {
  return (
    <FooterWrapper>
      &copy; Chuprakov Vadim, 2022. All Rights Reserved.
      <br />
      Contact me: <a href="https://t.me/cprkv">@cprkv</a>
    </FooterWrapper>
  );
}
