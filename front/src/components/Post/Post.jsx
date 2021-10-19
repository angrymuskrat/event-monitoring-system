import React from 'react'
import PropTypes from 'prop-types'

// styled
import Container from './Post.styled'

function Post({ post }) {
  return (
    <Container key={post.id}>
      <div className="post__header">
        <a href={post.profileLink} target="_blank" rel="noopener noreferrer">
          <div
            className="post__profile-pic"
            style={{
              background: `url(${post.profilePicUrl})`,
              backgroundPosition: 'center',
              backgroundRepeat: 'no-repeat',
              backgroundSize: 'cover',
            }}
          ></div>
        </a>

        <div className="post__profile-info">
          <a
            className="text text_bold text_post"
            href={post.profileLink}
            target="_blank"
            rel="noopener noreferrer"
          >
            {post.username}
          </a>

          <p className="text text_p2 text_location">{post.location}</p>
        </div>
        <a
          href={post.profileLink}
          target="_blank"
          rel="noopener noreferrer"
          className="text text_p2"
        >
          <button className="post__profile-button">View profile</button>
        </a>
      </div>
      <div className="post__picture">
        {/* <a href={post.postLink} target="_blank" rel="noopener noreferrer">
          <img src={post.photoUrl} alt={post.caption} />
        </a> */}
        <img src={post.photoUrl} alt={post.caption} />
      </div>
      <div className="post__likes">
        <p className="text text_post">
          ♡ {post.likes} {post.likes === 1 ? 'like' : 'likes'}
        </p>
      </div>
      <div className="post__description">
        <a
          href={post.profileLink}
          target="_blank"
          rel="noopener noreferrer"
          className="text text_bold text_post"
        >
          {post.username}
        </a>
        <p className="text text_post">{post.caption}</p>
      </div>
    </Container>
  )
}
Post.propTypes = {
  post: PropTypes.object.isRequired,
}

export default Post
