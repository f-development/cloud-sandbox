use test;
explain select `id`, max(`created_at`) from `test1` where `group` in ('1', '2', '3') group by `id`;