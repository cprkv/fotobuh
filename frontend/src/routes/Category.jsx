import React from "react";
import { useParams } from "react-router-dom";
import styled from "styled-components";
import { MainWrapper } from "../components/MainWrapper";
import { apiGetCategory, apiImageUrl } from "../api";

const SimpleGrid = styled.div`
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  align-items: center;
  padding: 0 8px;

`;

const GridImage = styled.img`
  max-height: 95vh;
  max-width: 100vw;  /* TODO! FIX HORIZONTAL SCROLL !!!! */
  margin: 1rem;
`

export function Category() {
  const params = useParams();

  if (!params.categoryId) {
    return <MainWrapper>not found</MainWrapper>;
  }

  const [category, setCategory] = React.useState([]);
  React.useEffect(() => {
    apiGetCategory(params.categoryId).then(setCategory);
  }, [params.categoryId]);

  return (
    <SimpleGrid>
      {category &&
        category.pictures &&
        category.pictures.length &&
        category.pictures.map(({ fileName, name, createdAt }, index) => (
          <GridImage src={apiImageUrl(fileName)} alt={name} key={index} />
        ))}
    </SimpleGrid>
  );
}
