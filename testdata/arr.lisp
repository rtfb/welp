(let arr (mk-array))
(append arr 3 4 5)
(let i 1)
(nth i (append arr 7))
(nth 2 (append (mk-array) 3 (+ 2 3) 7))
(fn car (arr) (nth 0 arr))
(car arr)
(fn +1 (arg) (+ arg 1))

(fn rest-impl (arr new-arr pos)
  (cond
    ((eq pos (len arr)) new-arr)
    (t
      (rest-impl arr (append new-arr (nth pos arr)) (+1 pos)))))

(fn rest (arr)
  (rest-impl arr (mk-array) 1))

(rest arr)
