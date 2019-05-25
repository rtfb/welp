(fn car (arr) (nth 0 arr))
(fn first (arr) (car arr))
(fn rest-impl (arr new-arr pos)
  (cond
    ((eq pos (len arr)) new-arr)
    (t
      (rest-impl arr (append new-arr (nth pos arr)) (+1 pos)))))
(fn rest (arr)
  (rest-impl arr (mk-array) 1))
(fn cdr (arr) (rest arr))
(fn +1 (arg) (+ arg 1))
(fn 1+ (arg) (+1 arg))
