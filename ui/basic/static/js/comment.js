const commentSubmit = document.getElementById('comment-submit');
const comment = document.getElementById('comment');
const commentAlert = document.getElementById('comment-alert');
const commentAlertHide = document.getElementById('comment-alert-hide');

commentAlertHide.addEventListener('click', () => {
  commentAlert.classList.add('hidden');
});

commentSubmit.addEventListener('click', async () => {
  commentAlert.classList.add('hidden');

  if (!(comment instanceof HTMLTextAreaElement)) {
    return;
  }

  if (!comment.value) {
    return;
  }

  const slug = commentSubmit.getAttribute('data-slug');

  if (!slug) {
    return;
  }

  try {
    const resp = await fetch(`/comments/${slug}/create`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
      body: `Message=${encodeURIComponent(comment.value)}&ParentId=null`,
    });

    if (resp.ok) {
      comment.value = '';

      commentAlert.classList.remove('hidden');

      return;
    }

    throw new Error();
  } catch (error) {
    // TODO
  }
});

const commentReplies = document.querySelectorAll('[data-comment="reply"]');

commentReplies.forEach((el) => {
  const parent = el.getAttribute('data-parent');
  const replySection = document.querySelector(
    `[data-comment="reply-section"][data-parent="${parent}"]`
  );

  el.addEventListener('click', () => {
    replySection.classList.toggle('hidden');
  });
});

const commentSubmits = document.querySelectorAll('[data-comment="submit"]');

commentSubmits.forEach((el) => {
  const parent = el.getAttribute('data-parent');
  const commentChild = document.querySelector(
    `[data-comment="comment"][data-parent="${parent}"]`
  );
  const commentAlertChild = document.querySelector(
    `[data-comment="alert"][data-parent="${parent}"]`
  );
  const commentAlertHideChild = document.querySelector(
    `[data-comment="alert-hide"][data-parent="${parent}"]`
  );

  commentAlertHideChild.addEventListener('click', () => {
    commentAlertChild.classList.add('hidden');
  });

  el.addEventListener('click', async () => {
    commentAlertChild.classList.add('hidden');

    if (!(commentChild instanceof HTMLTextAreaElement)) {
      return;
    }

    if (!commentChild.value) {
      return;
    }

    const slug = el.getAttribute('data-slug');

    if (!slug) {
      return;
    }

    try {
      const resp = await fetch(`/comments/${slug}/create`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: `Message=${encodeURIComponent(commentChild.value)}&ParentId=${parent}`,
      });

      if (resp.ok) {
        commentChild.value = '';

        commentAlertChild.classList.remove('hidden');

        return;
      }

      throw new Error();
    } catch (error) {
      // TODO
    }
  });
});
