import React, { useEffect, useState } from "react";
import PropTypes from "prop-types";

// styled
import Container from "./Post.styled";

import noPhoto from "../../assets/img/noPhoto.png";

// const checkImgSrc = (src) => {
//   const img = new Image();
//   img.onload = function() {
//     console.log(`valid src: ${src}`);
//   };
//   img.onerror = function() {
//     console.log(`unvalid src: ${src}`);
//   };
//   img.src = src;
// };

function Post({ post }) {
  const [isPhotoValid, setIsPhotoValid] = useState(noPhoto);

  useEffect(() => {
    // checkImgSrc(post.photoUrl);
    // если фото отгрузилось отображаем его если нет то ставим картинку
    const img = new Image();
    img.onload = () => {
      setIsPhotoValid(img.src);
      // console.log(`img`, img.src);
    };
    img.onerror = function() {
      setIsPhotoValid(noPhoto);
    };
    img.src = post.photoUrl;
  }, [post.photoUrl]);

  return (
    <Container key={post.id}>
      <div className="post__header">
        {/* <a href={post.profileLink} target="_blank" rel="noopener noreferrer">
          <div
            className="post__profile-pic"
            style={{
              background: `url(${post.profilePicUrl})`,
              backgroundPosition: "center",
              backgroundRepeat: "no-repeat",
              backgroundSize: "cover",
            }}
          ></div>
        </a> */}

        <div className="post__profile-info">
          <p className="text text_bold text_post">{post.username}</p>

          <a
            className="text text_p2"
            href={post.locationLink}
            target="_blank"
            rel="noopener noreferrer"
          >
            <span>location id: </span>
            <span className="text_location">{post.location}</span>
          </a>
        </div>
        <a
          href={post.postLink}
          target="_blank"
          rel="noopener noreferrer"
          className="text text_p2"
        >
          <button className="post__profile-button">View post</button>
        </a>
      </div>
      <div className="post__picture">
        {/* <a href={post.postLink} target="_blank" rel="noopener noreferrer">
          <img src={post.photoUrl} alt={post.caption} />
        </a> */}
        {/* <img src={post.photoUrl} alt={post.caption} /> */}
        <img src={isPhotoValid} alt={post.caption} />
        {/* {console.log(`post.photoUrl`, post.photoUrl)} */}
      </div>
      <div className="post__likes">
        <p className="text text_post">
          ♡ {post.likes} {post.likes === 1 ? "like" : "likes"} | {post.comments}{" "}
          {post.likes === 1 ? "comment" : "comments"}
        </p>
      </div>
      <div className="post__description">
        <p className="text text_bold text_post">Description</p>
        <p className="text text_post">{post.caption}</p>
      </div>
    </Container>
  );
}
Post.propTypes = {
  post: PropTypes.object.isRequired,
};

export default Post;
