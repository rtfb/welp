(fn fib (n)
  (cond
    ((eq n 1) 1)
    ((eq n 2) 1)
    (t (+ (fib (- n 1)) (fib (- n 2))))))
(fib 7)
