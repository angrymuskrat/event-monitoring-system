import styled from 'styled-components'
export default styled.div`
  display: flex;
  flex-direction: column;
  max-width: 520px;
  width: 100%;
  margin: 0 auto;

  .post__header {
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem;
  }
  .post__profile-pic {
    background-color: #fff;
    background-size: cover;
    width: 100%;
    height: 100%;
    width: 4rem;
    height: 4rem;
    border-radius: 50%;
  }
  .post__profile-info {
    flex-basis: 65%;
    text-overflow: ellipsis;
    overflow: hidden;
  }
  .post__profile-button {
    cursor: pointer;
    background-color: #3897f0;
    color: #fff;
    padding: 0.8rem;
    transition: all 0.3s;
    font-size: 1.3rem;
    &:hover {
      background-color: #1372cc;
    }
    @media (max-width: 375px) {
      font-size: 0.8rem;
    }
  }
  .post__likes {
    padding-left: 1.5rem;
    padding-right: 1.5rem;
  }
  .post__picture {
    background-color: #fff;
    margin-bottom: 1rem;
    & img {
      width: 100%;
      height: auto;
    }
  }
  .post__description {
    max-width: 52rem;
    padding-left: 1.5rem;
    padding-right: 1.5rem;
    padding-bottom: 2.5rem;
    a {
      display: inline-block;
      padding-bottom: 0.4rem;
    }
    .text {
      text-overflow: ellipsis;
      word-wrap: break-word;
    }
  }
`
